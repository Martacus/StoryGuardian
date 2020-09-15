package com.guardian.custodian;

import lombok.Data;
import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.util.List;
import java.util.Set;
import java.util.UUID;

@Data
@javax.persistence.Entity
public class Entity {

    @Id @GeneratedValue(generator="system-uuid")
    @GenericGenerator( name = "UUID", strategy = "org.hibernate.id.UUIDGenerator")
    private UUID id;

    @NotNull
    @Size(min = 1, max = 50)
    private String name;

    @Size(max = 160)
    private String description;

    @ManyToOne(optional = true)
    private Entity parent;
    
    @OneToMany
    private Set<Entity> children;
}
