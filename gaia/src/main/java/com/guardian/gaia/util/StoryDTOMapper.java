package com.guardian.gaia.util;

import org.modelmapper.ModelMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;

import java.text.ParseException;

@Component
public class StoryDTOMapper<T, J> {

    private final ModelMapper modelMapper;

    public StoryDTOMapper(ModelMapper modelMapper) {
        this.modelMapper = modelMapper;
    }

    public J convertToDto(T from, Class<J> to) {
        J postDto = modelMapper.map(from, to);
        return postDto;
    }

    public T convertToEntity(J from, Class<T> to) throws ParseException {
        T post = modelMapper.map(from, to);
        return post;
    }

}
