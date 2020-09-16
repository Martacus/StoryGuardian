package com.guardian.gaia.dto;

import com.guardian.custodian.Entity;
import lombok.Data;

import javax.persistence.ManyToOne;
import javax.persistence.OneToMany;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.util.Set;
import java.util.UUID;

@Data
public class EntityDTO {
    private UUID id;
    private String name;
    private String description;
    private Entity parent;
    private Set<Entity> children;
}
