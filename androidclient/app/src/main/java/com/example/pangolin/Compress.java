package com.example.pangolin;

import android.util.Log;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.nio.ByteBuffer;
import java.util.zip.GZIPInputStream;
import java.util.zip.GZIPOutputStream;

public class Compress {
    static final int MAX_PACKET_SIZE = PangolinVpnService.MAX_PACKET_SIZE;

    public static byte[] compress(byte[] data) {
        try {
            ByteArrayOutputStream out = new ByteArrayOutputStream();
            GZIPOutputStream gzip = new GZIPOutputStream(out);
            gzip.write(data);
            gzip.flush();
            gzip.close();
            return out.toByteArray();

        }catch(Exception e){
            Log.e("Compress", e.toString());
        }
        return null;
    }

    public static byte[] uncompress(byte[] data){
        try {
            ByteArrayOutputStream out = new ByteArrayOutputStream();
            ByteArrayInputStream in = new ByteArrayInputStream(data);
            GZIPInputStream gunzip = new GZIPInputStream(in);
            byte[] buffer = new byte[MAX_PACKET_SIZE];
            int n;
            while ((n = gunzip.read(buffer)) >= 0) {
                out.write(buffer, 0, n);
            }
            out.flush();
            return out.toByteArray();

        }catch(Exception e){
            Log.e("uncompress: ", e.toString());

        }
        return null;
    }
}
