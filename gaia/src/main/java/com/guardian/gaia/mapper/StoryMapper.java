package com.guardian.gaia.mapper;

import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.util.StoryDTOMapper;
import org.modelmapper.ModelMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class StoryMapper extends StoryDTOMapper<Story, StoryDTO> {

    public StoryMapper(ModelMapper modelMapper) {
        super(modelMapper, Story.class, StoryDTO.class);
    }




}
