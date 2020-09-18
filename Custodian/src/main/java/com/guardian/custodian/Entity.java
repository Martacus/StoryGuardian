package com.guardian.custodian;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.ToString;
import org.hibernate.annotations.GenericGenerator;

import javax.persistence.*;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.io.Serializable;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.UUID;

@Data
@javax.persistence.Entity
@Table
@EqualsAndHashCode(exclude="parent")
public class Entity implements Serializable {

    @Id @GeneratedValue(generator="system-uuid")
    @GenericGenerator( name = "UUID", strategy = "org.hibernate.id.UUIDGenerator")
    private UUID id;

    @NotNull
    @Size(min = 1, max = 50)
    private String name;

    @Size(max = 160)
    private String description;

    @ManyToOne(fetch = FetchType.LAZY)
    private Story story;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "parent_id")
    private Entity parent;
    
    @OneToMany(cascade=CascadeType.ALL, fetch = FetchType.LAZY, mappedBy="parent")
    private Set<Entity> children = new HashSet<>();

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public void linkToParent(Entity parent){
        this.parent = parent;
    }

    public void addChild(Entity entity){
        this.children.add(entity);
    }

    @JsonIgnore
    public Set<Entity> getChildren() {
        return children;
    }
}
