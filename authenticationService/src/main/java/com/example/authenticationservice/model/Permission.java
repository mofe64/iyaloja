package com.example.authenticationservice.model;

import com.fasterxml.jackson.annotation.JsonAlias;
import lombok.Data;

@Data
public class Permission {
    @JsonAlias({"_id"})
    private String id;
    private String name;
}
