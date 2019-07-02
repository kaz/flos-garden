<template>
  <main>
    <header>FLOS ❀ Hortus</header>
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
body, html {
  padding: 0;
  margin: 0;
  font-family: monospace;
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
  padding: 1em 2em 0.7em 2em;
  cursor: pointer;
}
nav > div.now {
  border-bottom: 0.3em solid orange;
}

section {
  padding: 1em 2em;
}
</style>
