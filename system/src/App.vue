<template>
  <nav class="navbar navbar-light guardian-nb fixed-top">

    </nav>
  <div class="d-flex"  id="wrapper">
    

    <nav id="sidebar" class="sticky-top ">
        <ul class="list-unstyled components">
            <!-- <li class="active">
                <a href="#homeSubmenu" data-toggle="collapse" aria-expanded="false" class="dropdown-toggle">Home</a>
                <ul class="collapse list-unstyled" id="homeSubmenu">
                    <li>
                        <a href="#">Home 1</a>
                    </li>
                    <li>
                        <a href="#">Home 2</a>
                    </li>
                    <li>
                        <a href="#">Home 3</a>
                    </li>
                </ul>
            </li> -->
            <li>
                <router-link to="/portal">About</router-link>
            </li>
        </ul>
    </nav>

    <div class="route-view">
      <div class="container">
        <router-view/>
      </div>
    </div>
  </div>
  
</template>

<style>
@import './css/site.css';
</style>
