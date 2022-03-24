package com.example.authenticationservice.service;

import com.example.authenticationservice.config.FeignConfig;
import com.example.authenticationservice.model.LoginRequest;
import com.example.authenticationservice.model.LoginResponse;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestHeader;

@FeignClient(value = "userService", url = "http://localhost:4000/api/v1", configuration = FeignConfig.class)
public interface UserServiceClient {
    @PostMapping(value = "/auth/login", consumes = "application/json", produces = "application/json")
    LoginResponse login(@RequestBody LoginRequest request);
}
