package com.example.authenticationservice.model;

import lombok.Data;

@Data
public class LoginResponse {
    private String status;
    private String token;
    private com.example.authenticationservice.model.Data data;
}
