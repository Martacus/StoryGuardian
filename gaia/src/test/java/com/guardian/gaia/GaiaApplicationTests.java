package com.guardian.gaia;

import com.guardian.custodian.Entity;
import com.guardian.gaia.repository.EntityRepository;
import com.guardian.gaia.service.EntityService;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit.jupiter.SpringExtension;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.context.junit4.SpringRunner;
import org.springframework.test.context.support.AnnotationConfigContextLoader;
import org.springframework.transaction.annotation.Transactional;

import javax.annotation.Resource;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.assertEquals;

@RunWith(SpringRunner.class)
@DataJpaTest
class GaiaApplicationTests {

//	@Autowired
//	private EntityService service;

	@Autowired
	private EntityRepository repository;


	@Test
	public void givenStudent_whenSave_thenGetOk() {
		Entity entity = new Entity();
		entity.setId(UUID.randomUUID());

		Entity entity2 = new Entity();
		entity2.setId(UUID.randomUUID());

		entity2 = repository.save(entity2);
		entity = repository.save(entity);

		//entity.linkAsParent(entity2);
		//entity2.addChild(entity);

		assertEquals(entity,entity);
	}

}
