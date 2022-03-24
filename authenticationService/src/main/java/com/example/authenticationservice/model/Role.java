package com.example.authenticationservice.model;

import com.fasterxml.jackson.annotation.JsonAlias;
import lombok.Data;

import java.util.Set;

@Data
public class Role {
    @JsonAlias({"_id"})
    private String id;
    private String name;
    private Set<Permission> permissions;
}
