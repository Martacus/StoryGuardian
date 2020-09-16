package com.guardian.gaia.mapper;

import com.guardian.custodian.Entity;
import com.guardian.gaia.dto.EntityDTO;
import com.guardian.gaia.util.DTOMapper;
import org.modelmapper.ModelMapper;
import org.springframework.stereotype.Component;

@Component
public class EntityMapper extends DTOMapper<Entity, EntityDTO> {
    public EntityMapper(ModelMapper modelMapper) {
        super(modelMapper, Entity.class, EntityDTO.class);
    }
}
