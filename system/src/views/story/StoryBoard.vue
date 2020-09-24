<template>
  <div class="row">
      <div class="col-md-8">
        <h2>Stories</h2>
      </div>

      <div class="col-md-4">
         <router-link to="/portal/story/add" tag="div">
           <a style="float: right;">
            <button type="button" class="btn guardian-btn ">
              <svg width="2em" height="2em" viewBox="0 0 16 16" class="bi bi-plus icon text-success" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4z"/>
              </svg>
            </button>
          </a> 
         </router-link>
      </div>
    </div>


    <div class="row">

      <div v-for="(item, i) in stories" :key="i" class="col-md-3">
        <router-link :to="item.name" tag="div">
          <div class="card" style="width: 100%">
            <div class="card-body" style="text-align: center;">
              <h1>{{item.name}}</h1>
            </div>
         </div>
         </router-link>
      </div>
    </div>
</template>

<script>
import instance from '../../controller/StoryController';

export default {
  name: 'StoryBoard',
  data() {
    return {
      stories: [
        {
          name: "story1"
        },
        {
          name: "story2"
        }
      ]
    }
  },
  created() {
    instance.get('/story/user/89787').then(response => {this.stories = response.data});
  }
  
} 
</script>

<style>

.icon{
  vertical-align: text-bottom;
}


</style>