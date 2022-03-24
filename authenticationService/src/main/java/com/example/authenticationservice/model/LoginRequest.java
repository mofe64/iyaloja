package com.example.authenticationservice.model;

import lombok.Data;

@Data
public class LoginRequest {
    private String email;
    private String password;
}
