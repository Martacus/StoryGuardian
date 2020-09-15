package com.guardian.gaia.mapper;

import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.util.DTOMapper;
import org.modelmapper.ModelMapper;
import org.springframework.stereotype.Component;

@Component
public class StoryMapper extends DTOMapper<Story, StoryDTO> {

    public StoryMapper(ModelMapper modelMapper) {
        super(modelMapper, Story.class, StoryDTO.class);
    }




}
