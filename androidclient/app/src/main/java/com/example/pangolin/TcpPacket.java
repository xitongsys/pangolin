package com.example.pangolin;

import java.io.InputStream;
import java.io.OutputStream;
import java.util.ArrayList;
import android.util.Log;

import static java.lang.Math.min;

public class TcpPacket {
    final static int MAXPACKETSIZE = 1024*1024;

    private static void writeEnd(OutputStream outputStream) throws Exception {
        byte[] bs = {0};
        outputStream.write(bs);
        outputStream.flush();
    }

    public static void write(byte[] data, OutputStream outputStream) throws Exception{
        int ln = data.length;
        int left = ln;
        while(left > 0){
            int wc = min(255, left);
            byte[] bs = {(byte)(wc)};
            outputStream.write(bs);
            bs = new byte[wc];
            for(int i=0; i<wc; i++){
                bs[i] = data[ln - left + i];
            }
            outputStream.write(bs);
            left -= wc;
        }
        writeEnd(outputStream);
    }

    public static int read(byte[] data, InputStream inputStream) throws Exception{
        byte[] bs;
        ArrayList<Byte> buf = new ArrayList<>();
        while(true){
            bs = new byte[1];
            inputStream.read(bs);
            int ln = (bs[0] & 0xFF);
            if(ln <= 0) break;

            int left = ln;
            while(left > 0){
                bs = new byte[left];
                int n = inputStream.read(bs);
                for(int i=0; i<n; i++){
                    buf.add(bs[i]);
                }
                left -= n;
            }
        }
        int lp = buf.size();
        for(int i=0; i<lp; i++){
            data[i] = buf.get(i);
        }
        return lp;
    }
}
