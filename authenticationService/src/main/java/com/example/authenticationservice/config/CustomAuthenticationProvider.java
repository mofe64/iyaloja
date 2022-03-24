package com.example.authenticationservice.config;

import com.example.authenticationservice.model.LoginRequest;
import com.example.authenticationservice.model.LoginResponse;
import com.example.authenticationservice.service.UserServiceClient;
import feign.FeignException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.authentication.AuthenticationServiceException;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.stereotype.Component;

import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;

@Slf4j
@Component
public class CustomAuthenticationProvider implements AuthenticationProvider {

    @Autowired
    UserServiceClient userServiceClient;

    @Override
    public Authentication authenticate(Authentication authentication) throws AuthenticationException {
        String name = authentication.getName();
        String password = authentication.getCredentials().toString();
        log.info("auth details --> {}, {}", name, password);
        try {
            LoginRequest request = new LoginRequest();
            request.setEmail(name);
            request.setPassword(password);
            LoginResponse response = userServiceClient.login(request);
            log.info("response --> {}", response);
            assert response != null;
            var username = response.getData().getUser().getEmail();
            var userPassword = response.getData().getUser().getPassword();
            var roles = response.getData().getUser().getRoles();
            Set<SimpleGrantedAuthority> authorities = new HashSet<>();
            roles.forEach(
                    role -> {
                        Set<SimpleGrantedAuthority> permissions = role.getPermissions().stream()
                                .map(permission -> new SimpleGrantedAuthority(permission.getName()))
                                .collect(Collectors.toSet());
                        permissions.add(new SimpleGrantedAuthority(role.getName()));
                        authorities.addAll(permissions);
                    }
            );
            return new UsernamePasswordAuthenticationToken(
                    username,
                    userPassword,
                    authorities
            );
        } catch (FeignException.Unauthorized e) {
            throw new BadCredentialsException("Bad Authentication credentials provided");
        } catch (FeignException e) {
            throw new AuthenticationServiceException("Could not authenticate");
        }
    }

    @Override
    public boolean supports(Class<?> authentication) {
        return authentication.equals(UsernamePasswordAuthenticationToken.class);
    }


}
