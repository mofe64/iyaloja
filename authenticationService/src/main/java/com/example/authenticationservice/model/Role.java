package com.example.authenticationservice.model;

import com.fasterxml.jackson.annotation.JsonAlias;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.Set;

@Data
public class Role {
    @JsonAlias({"_id"})
    private String id;
    private String name;
    private Set<Permission> permissions;
}
