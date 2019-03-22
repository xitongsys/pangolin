package com.example.pangolin;

import android.app.Activity;
import android.content.Intent;
import android.content.SharedPreferences;
import android.net.VpnService;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.SocketAddress;
import java.nio.ByteBuffer;
import java.nio.channels.DatagramChannel;

public class MainActivity extends AppCompatActivity {
    private Button btConn, btDisconn;
    private EditText editServer, editServerPort, editLocal, editDNS;
    private TextView viewInfo;
    private Thread sendThread;
    SharedPreferences preferences;
    SharedPreferences.Editor preEditor;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        btConn = findViewById(R.id.connButton);
        btDisconn = findViewById(R.id.disconnButton);
        editServer = findViewById(R.id.serverAddrEdit);
        editServerPort = findViewById(R.id.serverPortEdit);
        editLocal = findViewById(R.id.localAddrEdit);
        viewInfo = findViewById(R.id.infoTextView);
        editDNS = findViewById(R.id.dnsEdit);

        sendThread = null;

        preferences = getPreferences(Activity.MODE_PRIVATE);
        preEditor = preferences.edit();

        editServer.setText(preferences.getString("serverIP", "192.168.0.1"));
        editServerPort.setText(preferences.getString("serverPort", "12345"));
        editLocal.setText(preferences.getString("localIP", "10.0.0.3/8"));
        editDNS.setText(preferences.getString("dns", "8.8.8.8"));

        btDisconn.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                viewInfo.setText("Disconnect");
                Intent intent = new Intent();
                intent.setClass(MainActivity.this, PangolinVpnService.class);
                intent.setAction("disconnect");
                startService(intent);
            }
        });

        btConn.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                viewInfo.setText("Connect");
                Intent intent = VpnService.prepare(MainActivity.this);
                if(intent != null){
                    startActivityForResult(intent, 0);
                }else{
                    onActivityResult(0, RESULT_OK, null);
                }
            }
        });

    }

    @Override
    protected void onActivityResult(int request, int result, Intent data) {
        if (result == RESULT_OK) {
            Intent intent = new Intent(this, PangolinVpnService.class);

            String serverIP = editServer.getText().toString();
            String serverPort = editServerPort.getText().toString();
            String localIp = editLocal.getText().toString();
            String dns = editDNS.getText().toString();

            intent.setAction("connect");
            intent.putExtra("serverIP", serverIP);
            intent.putExtra("serverPort", Integer.parseInt(serverPort));
            intent.putExtra("localIP", localIp);
            intent.putExtra("dns", dns);
            startService(intent);

            preEditor.putString("serverIP", serverIP);
            preEditor.putString("serverPort", serverPort);
            preEditor.putString("localIP", localIp);
            preEditor.putString("dns", dns);
            preEditor.commit();

        }
    }

}
