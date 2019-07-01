package com.example.pangolin;

import android.app.Notification;
import android.app.PendingIntent;
import android.content.Intent;
import android.net.VpnService;
import android.os.Bundle;
import android.os.IBinder;
import android.os.ParcelFileDescriptor;
import android.util.Log;

import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.net.Socket;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.DatagramChannel;
import java.util.Arrays;

public class PangolinVpnService extends VpnService {
    final static String ACTION_DISCONNECT = "disconnect";
    final static String ACTION_CONNECT = "connect";
    final static int MAX_PACKET_SIZE = 65536;
    static String serverIP, localIP;
    static int localPrefixLength = 24;
    static int serverPort;
    static String dns;
    static String protocol = "tcp";
    static String token = "";
    Thread sendrecvThreadUdp;
    Thread sendThreadTcp, recvThreadTcp;
    Socket tcpSocket;
    ParcelFileDescriptor localTunnel;
    private PendingIntent pendingIntent;
    private Encryption encryption;

    public PangolinVpnService() {
    }

    @Override
    public void onCreate(){
        pendingIntent = PendingIntent.getActivity(this, 0, new Intent(this, MainActivity.class),
                PendingIntent.FLAG_CANCEL_CURRENT);
    }

    @Override
    public IBinder onBind(Intent intent) {
        return null;
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        Log.i("onStartCommand", "start: " + intent.getAction());
        try {
            if (intent != null && ACTION_DISCONNECT.equals(intent.getAction())) {
                disconnect();
                return START_NOT_STICKY;
            } else {
                Bundle ex = intent.getExtras();
                serverIP = ex.getString("serverIP");
                serverPort = ex.getInt("serverPort");
                protocol = ex.getString("protocol");
                String[] localAddrs = ex.getString("localIP").split("/");
                if(localAddrs.length>=1){
                    localIP = localAddrs[0];
                }
                if(localAddrs.length>=2){
                    localPrefixLength = Integer.parseInt(localAddrs[1]);
                }

                dns = ex.getString("dns");
                token = ex.getString("token");
                encryption = new Encryption(token);

                Notification.Builder builder = new Notification.Builder(this);
                builder.setContentIntent(pendingIntent)
                        .setSmallIcon(R.mipmap.ic_launcher)
                        .setContentTitle("Pangolin")
                        .setContentText("<Server>" + serverIP + ":" + serverPort)
                        .setWhen(System.currentTimeMillis());
                Notification notification = builder.build();
                startForeground(1, notification);
                connect();
            }
        }catch (Exception e){
            Log.e("onStartCommmand", e.toString());
        }
        return START_STICKY;
    }

    private void initUdpThread() {
        sendrecvThreadUdp = new Thread() {
            @Override
            public void run() {
                try {
                    final DatagramChannel udp = DatagramChannel.open();

                    SocketAddress serverAdd = new InetSocketAddress(serverIP, serverPort);
                    udp.connect(serverAdd);
                    udp.configureBlocking(false);
                    PangolinVpnService.this.protect(udp.socket());

                    VpnService.Builder builder = PangolinVpnService.this.new Builder();
                    builder.setMtu(1500)
                            .addAddress(localIP, localPrefixLength)
                            .addRoute("0.0.0.0", 0)
                            .addDnsServer(dns)
                            .setSession("Pangolin")
                            .setConfigureIntent(null);
                    localTunnel = builder.establish();


                    FileInputStream in = new FileInputStream(localTunnel.getFileDescriptor());
                    FileOutputStream out = new FileOutputStream(localTunnel.getFileDescriptor());

                    while(!isInterrupted()){
                        try {
                            byte[] buf = new byte[MAX_PACKET_SIZE];
                            int ln = in.read(buf);
                            if (ln > 0) {
                                byte[] cmpbuf = Compress.compress(Arrays.copyOfRange(buf, 0, ln));
                                ByteBuffer bf = ByteBuffer.wrap(cmpbuf);
                                udp.write(bf);
                            }

                            ByteBuffer bf = ByteBuffer.allocate(MAX_PACKET_SIZE);
                            ln = udp.read(bf);
                            if (ln > 0) {
                                bf.limit(ln); bf.rewind();
                                buf = new byte[ln];
                                bf.get(buf);
                                byte[] undata = Compress.uncompress(buf);
                                out.write(undata);
                            }

                        }catch(Exception e){
                            Log.e("send/rec", e.toString());
                        }
                    }

                }catch(Exception e){
                    Log.e("send/recv", e.toString());
                }
            }
        };
    }

    private void initTcpThread() {
        sendThreadTcp = new Thread(){
            @Override
            public void run(){
                try{
                    FileInputStream in = new FileInputStream(localTunnel.getFileDescriptor());
                    OutputStream out = tcpSocket.getOutputStream();
                    TcpPacket.write(Compress.compress(token.getBytes()), out);

                    while(!isInterrupted()){
                        try {
                            byte[] buf = new byte[MAX_PACKET_SIZE];
                            int ln = in.read(buf);
                            if (ln > 0) {
                                byte[] endata = encryption.encrypt(Arrays.copyOfRange(buf, 0, ln));
                                TcpPacket.write(Compress.compress(endata), out);
                            }

                        }catch(Exception e){
                            Log.e("sendThreadTcp", e.toString());
                        }
                    }

                }catch (Exception e){
                    Log.e("sendThreadTcp", e.toString());
                }
            }
        };

        recvThreadTcp = new Thread(){
            @Override
            public void run(){
                try{
                    FileOutputStream out = new FileOutputStream(localTunnel.getFileDescriptor());
                    InputStream in = tcpSocket.getInputStream();

                    while(!isInterrupted()){
                        try {
                            byte[] buf = new byte[MAX_PACKET_SIZE];
                            int ln = TcpPacket.read(buf, in);
                            if (ln > 0) {
                                byte[] undata = Compress.uncompress(Arrays.copyOfRange(buf, 0, ln));
                                out.write(encryption.decrypt(undata));
                            }

                        }catch(Exception e){
                            Log.e("recvThreadTcp", e.toString());
                        }
                    }

                }catch (Exception e){
                    Log.e("recvThreadTcp", e.toString());
                }
            }
        };
    }

    private void closeAll(){
        try {
            if (sendThreadTcp != null) {
                sendThreadTcp.interrupt();
                sendThreadTcp = null;
            }
            if (recvThreadTcp != null) {
                recvThreadTcp.interrupt();
                recvThreadTcp = null;
            }

            if (sendrecvThreadUdp != null) {
                sendrecvThreadUdp.interrupt();
                sendrecvThreadUdp = null;
            }

            if (localTunnel != null) {
                localTunnel.close();
                localTunnel = null;
            }
        }catch (Exception e){
            Log.e("closeAll", e.toString());
        }
    }

    private void disconnect(){
        Log.i("disconnect", "disconnecting...");
        try {
            closeAll();
            stopForeground(true);

        }catch(Exception e){
            Log.e("disconnect", e.toString());
        }
    }

    private void connect(){
        Log.i("connect", "connecting...");
        Log.i("vpn", serverIP + " " + serverPort + " " + localIP + " " + dns);
        try {
            closeAll();

            if(protocol.equals("udp")){
                initUdpThread();
                sendrecvThreadUdp.start();

            }else{
                initTcpThread();

                new Thread(){
                    @Override
                    public void run(){
                        try {
                            tcpSocket = new Socket(serverIP, serverPort);
                            tcpSocket.setKeepAlive(true);
                            PangolinVpnService.this.protect(tcpSocket);

                            VpnService.Builder builder = PangolinVpnService.this.new Builder();
                            builder.setMtu(1500)
                                    .addAddress(localIP, localPrefixLength)
                                    .addRoute("0.0.0.0", 0)
                                    .addDnsServer(dns)
                                    .setSession("Pangolin")
                                    .setConfigureIntent(null);
                            localTunnel = builder.establish();

                            sendThreadTcp.start();
                            recvThreadTcp.start();

                        }catch (Exception e){
                            Log.e("connect", e.toString());
                        }
                    }
                }.start();
            }
        }catch(Exception e){
            Log.e("vpn", e.toString());
        }
    }
}
