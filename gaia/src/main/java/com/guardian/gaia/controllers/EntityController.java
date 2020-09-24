package com.guardian.gaia.controllers;

import com.guardian.custodian.Entity;
import com.guardian.gaia.dto.EntityDTO;
import com.guardian.gaia.service.EntityService;
import com.guardian.gaia.util.DTOMapper;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;
import java.util.UUID;
import java.util.stream.Collectors;

@CrossOrigin(origins = "http://localhost:8000")
@RestController
@RequestMapping("/entity")
public class EntityController {

    private final EntityService entityService;
    private final DTOMapper<Entity, EntityDTO> entityMapper;

    public EntityController(EntityService entityService, DTOMapper<Entity, EntityDTO> entityMapper) {
        this.entityService = entityService;
        this.entityMapper = entityMapper;
    }


    //GET
    @GetMapping(value="/{uuid}")
    public EntityDTO getEntity(@PathVariable UUID uuid){
        Optional<Entity> entity = entityService.getById(uuid);
        return entity.map(entityMapper::convertToDto).orElse(null);
    }
    @GetMapping(value="/user/{uuid}")
    public List<EntityDTO> getAllEntities(@PathVariable String uuid){
        return entityService.getEntityList(uuid).stream()
                .map(entityMapper::convertToDto)
                .collect(Collectors.toList());
    }

    //POST
    @PostMapping("/")
    public EntityDTO PostEntity(@RequestBody EntityDTO entityDTO){
        Entity entity = entityMapper.convertToEntity(entityDTO);
        return entityMapper.convertToDto(entityService.saveEntity(entity));
    }

    @PostMapping("/{parentUUID}/link/{childUUID}")
    @ResponseStatus(HttpStatus.OK)
    public void LinkEntityToParentEntity(@PathVariable UUID parentUUID, @PathVariable UUID childUUID){
        entityService.linkEntity(parentUUID, childUUID);
    }


    //PUT
    @PutMapping("/")
    @ResponseStatus(HttpStatus.OK)
    public void PutEntity(@RequestBody EntityDTO entityDTO){
        Entity entity = entityMapper.convertToEntity(entityDTO);
        entityService.editEntity(entity);
    }

    //DELETE
    @DeleteMapping("/{uuid}")
    @ResponseStatus(HttpStatus.OK)
    public void DeleteEntity(@PathVariable UUID uuid){
        entityService.deleteEntity(uuid);
    }


}
