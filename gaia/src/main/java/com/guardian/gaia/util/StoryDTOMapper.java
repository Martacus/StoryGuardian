package com.guardian.gaia.util;

import org.modelmapper.ModelMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Component;

import java.text.ParseException;

public class StoryDTOMapper<T, J> {

    private final ModelMapper modelMapper;
    private final Class<T> modelClass;
    private final Class<J> dtoClass;

    public StoryDTOMapper(ModelMapper modelMapper, Class<T> modelClass, Class<J> dtoClass) {
        this.modelMapper = modelMapper;
        this.modelClass = modelClass;
        this.dtoClass = dtoClass;
    }

    public J convertToDto(T from) {
        J postDto = modelMapper.map(from, dtoClass);
        return postDto;
    }

    public T convertToEntity(J from) throws ParseException {
        T post = modelMapper.map(from, modelClass);
        return post;
    }

}
