package com.example.authenticationservice.model;

import com.fasterxml.jackson.annotation.JsonAlias;
import com.fasterxml.jackson.annotation.JsonIgnore;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;
import lombok.NoArgsConstructor;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.index.Indexed;
import org.springframework.data.mongodb.core.mapping.Document;

import java.util.HashSet;
import java.util.Set;


@Data

public class User {
    @JsonAlias({"_id"})
    private String id;
    private String firstName;
    private String lastName;
    private String email;
    private String password;
    private boolean active;
    private Set<Role> roles;

}
