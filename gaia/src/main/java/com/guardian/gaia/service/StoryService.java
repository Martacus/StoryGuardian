package com.guardian.gaia.service;

import com.guardian.custodian.Story;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

@Service
public class StoryService {

    public List<Story> getStoryList(String userid){
        List<Story> storyList = new ArrayList<>();
        storyList.add(new Story());
        return storyList;
    }

    public void saveStory(Story story) {

    }

    public void editStory(Story story) {

    }

}
