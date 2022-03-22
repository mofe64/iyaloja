package com.example.authenticationservice.config;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.oauth2.config.annotation.configurers.ClientDetailsServiceConfigurer;
import org.springframework.security.oauth2.config.annotation.web.configuration.AuthorizationServerConfigurerAdapter;
import org.springframework.security.oauth2.config.annotation.web.configurers.AuthorizationServerEndpointsConfigurer;
import org.springframework.security.oauth2.config.annotation.web.configurers.AuthorizationServerSecurityConfigurer;

import javax.annotation.Resource;

@Configuration
public class AuthorizationServerConfig extends AuthorizationServerConfigurerAdapter {
    @Autowired
    private AuthenticationManager authenticationManager;

    @Resource(name = "userService")
    private UserDetailsService userDetailsService;
    @Autowired
    private BCryptPasswordEncoder encoder;

    @Override
    public void configure(ClientDetailsServiceConfigurer clients) throws
            Exception {
        clients.inMemory()
                .withClient("iyaloja")
                .secret(encoder.encode("thisissecret"))
                .authorizedGrantTypes(
                        "password",
                        "client_credentials")
                .scopes("internal", "public");
    }

    @Override
    public void configure(AuthorizationServerEndpointsConfigurer endpoints) throws Exception {
        endpoints
                .authenticationManager(authenticationManager)
                .userDetailsService(userDetailsService);
    }

    @Override
    public void configure(AuthorizationServerSecurityConfigurer security) throws Exception {
        security.checkTokenAccess("permitAll()").tokenKeyAccess("permitAll()");

    }
}

