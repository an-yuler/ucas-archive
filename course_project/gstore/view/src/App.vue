<template>
  <div id="app">
    <el-container>
        <el-header>
          <el-menu
            :default="this.$route.path"
            background-color="#545c64"
            text-color="#fff"
            active-text-color="#ffd04b"
            mode="horizontal"
            router
          >
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/query1">查询1</el-menu-item>
            <el-menu-item index="/query2">查询2</el-menu-item>
            <el-menu-item index="/query3">查询3</el-menu-item>
          </el-menu>
        </el-header>
        <el-main>
          <!-- <el-header>
            <div id="nav">
              <router-link to="/">Home</router-link> |
              <router-link to="/about">About</router-link> |
              <router-link to="/query1">Q1</router-link> |
              <router-link to="/query2">Q2</router-link> |
              <router-link to="/query3">Q3</router-link>
            </div>
          </el-header> -->
          <router-view />
        </el-main>
    </el-container>
  </div>
</template>

<style>
#app {
  width: 70%;
  position: relative;                       
  left: 15%;     

  font-family: Avenir, Helvetica, Arial, sans-serif;
  /* -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale; */
  /* color: #2c3e50; */
}


#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
