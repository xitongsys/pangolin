package com.example.pangolin;

import android.content.Intent;
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

            intent.setAction("connect");
            intent.putExtra("serverIP", editServer.getText().toString());
            intent.putExtra("serverPort", Integer.parseInt(editServerPort.getText().toString()));
            intent.putExtra("localIP", editLocal.getText().toString());
            intent.putExtra("dns", editDNS.getText().toString());
            startService(intent);
        }
    }

}
