package com.guardian.custodian;

import lombok.Data;
import org.hibernate.annotations.GenericGenerator;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.util.UUID;

@Data
@Entity
public class Story {

    @Id @GeneratedValue(generator="system-uuid")
    @GenericGenerator( name = "UUID", strategy = "org.hibernate.id.UUIDGenerator")
    private UUID id;
    
    @NotNull
    @Size(min = 1, max = 50)
    private String name;

    @Size(max = 160)
    private String description;

    public UUID getId() {
        return id;
    }
}
