<template>
  <div class="row">
    <div class="col-md-12">
      <h2>Add Story</h2>
    </div>
  </div>  

  <div class="row">
    <div class="col-md-12">
      <form id="add-story-form" @submit="submitForm" style="width: 100%">
        <div class="row">
          <div class="col-md-6">
            <label>Story Name</label>
            <input v-model="storyname" class="form-control">
          </div>
          <div class="col-md-6">
            <label>Story Desc</label>
            <input v-model="storydesc" class="form-control">
          </div>
        </div>
        <div class="row">
          <div class="col-md-1">
            <button class="btn guardian-btn">Submit</button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import instance from '../../controller/StoryController';

export default {
  name: 'AddStory',
  data(){
    return {
      storyname: null,
      storydesc: null
    }
  },
  methods: {
    submitForm: function (){
      instance.post('/story',{
        name: this.storyname,
        description: this.storydesc
      });
    }
  }
}
</script>

<style>

</style>