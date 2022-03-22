package com.example.authenticationservice.service;

import com.example.authenticationservice.controller.RegRequest;

public interface UserService {
    public boolean register(RegRequest request);
}
