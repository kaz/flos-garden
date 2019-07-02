<template>
  <section id="backup">
    <h1>File Snapshots</h1>
    <table>
      <thead>
        <tr>
          <th>Created</th>
          <th>Host</th>
          <th>File</th>
          <th>Download</th>
        </tr>
      </thead>
      <tbody>
        <tr :key="ss.created" v-for="ss in snapshots">
          <td>{{ ss.created }}</td>
          <td>{{ ss.host }}</td>
          <td>{{ ss.file }}</td>
          <td><a :href="`/api/database/blob/${ss.host}/${ss.id}`">download</a></td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<style>
#backup table {
  border-collapse: collapse;
}
#backup th,
#backup td {
  padding: .5em 2em;
  border: 1px solid #999;
}
</style>

<script>
"use strict";

export default {
  data() {
    return {
      snapshots: [],
    }
  },
  created() {
    this.loadMonitors();
  },
  beforeDestroy() {
    clearInterval(this.timer);
  },
  methods: {
    async loadMonitors() {
      const resp = await fetch(`/api/database/query`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify("SELECT BINARY created, host, series, remote_id FROM bookshelf_data_archive ORDER BY created DESC"),
      });
      if(resp.status != 200){
        return alert(await resp.text());
      }

      this.snapshots = (await resp.json()).rows.map(([created, host, file, id]) => ({created, host, file, id}));
    }
  }
}
</script>
