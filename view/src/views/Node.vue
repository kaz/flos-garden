<template>
  <section>
    <h1>New</h1>
    <input type="text">
    <button>add</button>

    <h1>Bastion</h1>
    <table>
      <thead>
        <tr>
          <th>Hostname</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>

      </tbody>
    </table>

    <h1>Nodes</h1>
    <table>
      <thead>
        <tr>
          <th>Hostname</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>

      </tbody>
    </table>
  </section>
</template>

<style>
section {
  padding: 1em 2em;
}
table {
  border-collapse: collapse;
}
th, td {
  padding: .5em 2em;
  border: 1px solid #999;
}
</style>

<script>
export default {
}
</script>
