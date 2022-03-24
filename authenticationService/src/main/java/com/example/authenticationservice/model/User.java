package com.example.authenticationservice.model;

import com.fasterxml.jackson.annotation.JsonAlias;
import lombok.Data;
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
