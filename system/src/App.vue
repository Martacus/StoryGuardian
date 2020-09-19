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
                <router-link to="/portal/about">About</router-link>
            </li>
        </ul>
    </nav>

    <div class="route-view">
      <router-view/>
    </div>
  </div>
  
</template>

<style>
@import "https://fonts.googleapis.com/css?family=Poppins:300,400,500,600,700";

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  background-color: #f6f7f7;
}

#nav {
  padding: 30px;
}

#nav a {
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}

.guardian-nb{
  background-color: #1565a0;
  height: 47px;
  position: fixed;
  width: 100%;
}

.route-view{
  padding-top: 50px;
  padding-left: 50px;
  padding-right: 50px;
}

 .wrapper {
    display: flex;
    align-items: stretch;
}

#sidebar {
    min-width: 272px;
    max-width: 272px;
    min-height: 100vh;
}

#sidebar.active {
    margin-left: -250px;
}

a[data-toggle="collapse"] {
    position: relative;
}

.dropdown-toggle::after {
    display: block;
    position: absolute;
    top: 50%;
    right: 20px;
    transform: translateY(-50%);
}

body {
    font-family: 'Poppins', sans-serif;
    background: #fafafa;
}

a, a:hover, a:focus {
    color: inherit;
    text-decoration: none;
    transition: all 0.3s;
}

#sidebar {
    /* don't forget to add all the previously mentioned styles here too */
    background: #ffffff;
    color: #000000;
    transition: all 0.3s;
    padding-top: 50px;
    height: 100%;
    border-right: 1px solid #dcdcde;
}

#sidebar .sidebar-header {
    padding: 20px;
}

#sidebar ul.components {
    padding: 20px 0;
}

#sidebar ul p {
    color: #000000;
    padding: 10px;
}

#sidebar ul li a {
    padding: 10px;
    font-size: 1.1em;
    display: block;
}
#sidebar ul li a:hover {
    color: #000000;
    background: #fff;
}

#sidebar ul li.active > a, a[aria-expanded="true"] {
    color: #000000;
}
ul ul a {
    font-size: 0.9em !important;
    padding-left: 30px !important;
}

@media (max-width: 768px) {
    #sidebar {
        margin-left: -250px;
    }
    #sidebar.active {
        margin-left: 0;
    }
}

</style>
