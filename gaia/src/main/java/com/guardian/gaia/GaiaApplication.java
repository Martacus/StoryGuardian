package com.guardian.gaia;

import com.guardian.custodian.Story;
import com.guardian.gaia.dto.StoryDTO;
import com.guardian.gaia.util.StoryDTOMapper;
import org.modelmapper.ModelMapper;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@EntityScan(basePackages = "com.guardian.custodian")
public class GaiaApplication {

	public static void main(String[] args) {
		SpringApplication.run(GaiaApplication.class, args);
	}

	@Bean
	public ModelMapper modelMapper() {
		return new ModelMapper();
	}

}
