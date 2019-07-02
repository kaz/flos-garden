<template>
  <main>
    <header>FLOS ❀ ADMIN PANEL</header>
    <nav>
      <router-link tag="div" active-class="now" to="/node">Node</router-link>
      <router-link tag="div" active-class="now" to="/cnc">C&amp;C</router-link>
      <router-link tag="div" active-class="now" to="/monitor">Monitor</router-link>
      <router-link tag="div" active-class="now" to="/audit">Audit</router-link>
      <router-link tag="div" active-class="now" to="/backup">Backup</router-link>
    </nav>
    <router-view/>
  </main>
</template>

<style>
* {
  font-size: 14px;
  font-family: "Courier Prime", monospace;
}
body, html {
  padding: 0;
  margin: 0;
}
header {
  padding: 1em 2em;
  background-color: #111;
  color: #FFF;
}
nav {
  display: flex;
  width: 100vw;
  background-color: #333;
  color: #FFF;
}
nav > div {
  cursor: pointer;
  padding: 0.7em 2em 0.4em 2em;
  border-bottom: 0.3em solid transparent;
}
nav > div.now {
  border-bottom: 0.3em solid orange;
}

section {
  padding: 1em 2em;
}
</style>
