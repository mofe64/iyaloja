package com.example.authenticationservice.util;

import org.apache.commons.codec.binary.Base64;
import org.springframework.security.oauth2.common.util.SerializationUtils;
import org.springframework.security.oauth2.provider.OAuth2Authentication;

public class SerializableObjectConverter {

    public static String serialize(OAuth2Authentication object) {
        byte[] bytes = SerializationUtils.serialize(object);
        return Base64.encodeBase64String(bytes);
    }

    public static OAuth2Authentication deserialize(String encodedObject) {
        byte[] bytes = Base64.decodeBase64(encodedObject);
        return (OAuth2Authentication) SerializationUtils.deserialize(bytes);
    }
}
