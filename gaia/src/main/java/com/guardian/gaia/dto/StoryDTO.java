package com.guardian.gaia.dto;

import lombok.Data;

import java.util.UUID;

@Data
public class StoryDTO{
    private UUID id;
    private String name;
    private String description;
}
