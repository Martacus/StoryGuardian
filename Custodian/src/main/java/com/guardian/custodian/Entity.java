package com.guardian.custodian;

import lombok.Data;
import lombok.ToString;
import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.util.HashSet;
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

    @ManyToOne(fetch = FetchType.LAZY, optional = true)
    @ToString.Exclude
    private Entity parent;
    
    @OneToMany(fetch=FetchType.LAZY, cascade = CascadeType.ALL, mappedBy="parent")
    @ToString.Exclude
    private Set<Entity> children = new HashSet<>();

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public void linkAsParent(Entity entity){
        this.parent = entity;
        this.parent.addChild(entity);
    }

    public void addChild(Entity entity){
        this.children.add(entity);
    }

}
