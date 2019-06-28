package com.example.pangolin;

import java.security.GeneralSecurityException;
import java.util.Arrays;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;

public class Encryption {
    private IvParameterSpec ivSpec;
    private SecretKeySpec keySpec;

    public Encryption(String key){
        byte[] keyBytes = getKeyBytes(key);
        byte[] buf = new byte[16];
        System.arraycopy(keyBytes, 0, buf, 0, keyBytes.length > buf.length ? keyBytes.length : buf.length);
        ivSpec = new IvParameterSpec(keyBytes);
        keySpec = new SecretKeySpec(buf, "AES");
    }

    public byte[] encrypt(byte[] origData) throws GeneralSecurityException {
        Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
        cipher.init(Cipher.ENCRYPT_MODE, keySpec, ivSpec);
        return cipher.doFinal(origData);

    }

    public byte[] decrypt(byte[] crypted) throws GeneralSecurityException {
        Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
        cipher.init(Cipher.DECRYPT_MODE, keySpec, ivSpec);
        return cipher.doFinal(crypted);
    }

    private static byte[] getKeyBytes(String key) {
        byte[] bytes = key.getBytes();
        return bytes.length == 16 ? bytes : Arrays.copyOf(bytes, 16);
    }
}
