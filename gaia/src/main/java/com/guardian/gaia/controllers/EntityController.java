package com.guardian.gaia.controllers;

import com.guardian.custodian.Entity;
import com.guardian.gaia.dto.EntityDTO;
import com.guardian.gaia.service.EntityService;
import com.guardian.gaia.util.DTOMapper;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@RestController
@CrossOrigin(origins = "*")
@RequestMapping("/entity")
public class EntityController {

    private final EntityService entityService;
    private final DTOMapper<Entity, EntityDTO> entityMapper;

    public EntityController(EntityService entityService, DTOMapper<Entity, EntityDTO> entityMapper) {
        this.entityService = entityService;
        this.entityMapper = entityMapper;
    }


    //GET
    @GetMapping(value="/{uuid}", produces = {"application/json"})
    public ResponseEntity<EntityDTO> getEntity(@PathVariable UUID uuid){
        Optional<Entity> entity = entityService.getById(uuid);
        if(entity.isPresent()){
            return ResponseEntity.ok(entity.map(entityMapper::convertToDto).get());
        } else {
            return ResponseEntity.notFound().build();
        }
    }
    @GetMapping(value="/user/{uuid}",produces = {"application/json"})
    public ResponseEntity<List<EntityDTO>> getAllEntities(@PathVariable String uuid){
        return ResponseEntity.ok(
            entityService.getEntityList(uuid)
                .stream()
                .map(entityMapper::convertToDto)
                .collect(Collectors.toList())
        );
    }

    //POST
    @PostMapping(value="/", produces = {"application/json"})
    public ResponseEntity<EntityDTO> PostEntity(@RequestBody EntityDTO entityDTO){
        Entity entity = entityMapper.convertToEntity(entityDTO);
        return ResponseEntity.ok(entityMapper.convertToDto(entityService.saveEntity(entity)));
    }

    @PostMapping(value="/{parentUUID}/link/{childUUID}")
    @ResponseStatus(HttpStatus.OK)
    public void LinkEntityToParentEntity(@PathVariable UUID parentUUID, @PathVariable UUID childUUID){
        entityService.linkEntity(parentUUID, childUUID);
    }


    //PUT
    @PutMapping(value="/")
    @ResponseStatus(HttpStatus.OK)
    public void PutEntity(@RequestBody EntityDTO entityDTO){
        Entity entity = entityMapper.convertToEntity(entityDTO);
        entityService.editEntity(entity);
    }

    //DELETE
    @DeleteMapping(value="/{uuid}")
    @ResponseStatus(HttpStatus.OK)
    public void DeleteEntity(@PathVariable UUID uuid){
        entityService.deleteEntity(uuid);
    }


}
