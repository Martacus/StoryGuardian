package com.guardian.custodian;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.Id;
import java.util.UUID;

@Data
@Entity
public class Story {

    @Id
    public UUID id;
    public String name;
    public String description;

}
