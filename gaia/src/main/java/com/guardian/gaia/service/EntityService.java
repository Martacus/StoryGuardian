package com.guardian.gaia.service;

import com.guardian.custodian.Entity;
import com.guardian.gaia.repository.EntityRepository;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Component
public class EntityService {

    private final EntityRepository entityRepository;

    public EntityService(EntityRepository entityRepository) {
        this.entityRepository = entityRepository;
    }

    public List<Entity> getEntityList(String userid){
        return entityRepository.findAll();
    }

    public void saveEntity(Entity entity) {
        entityRepository.save(entity);
    }

    public void editEntity(Entity entity) {
        Optional<Entity> entity2 = entityRepository.findById(entity.getId());
        if(entity2.isPresent()){
            entityRepository.save(entity);
        }
    }

    public Optional<Entity> getById(UUID id) {
        return entityRepository.findById(id);
    }

    public void deleteEntity(UUID id) {
        entityRepository.findById(id).ifPresent(entityRepository::delete);
    }

    public void linkEntity(UUID parentUUID, UUID childUUID) {
        Optional<Entity> parent = entityRepository.findById(parentUUID);
        Optional<Entity> child = entityRepository.findById(childUUID);

        if(parent.isPresent() && child.isPresent()){
            child.get().linkAsParent(parent.get());
            //entityRepository.save(child.get());
            entityRepository.save(parent.get());
        }
    }
}
