<template>
  <section id="backup">
    <h1>Query</h1>
    Host: <input type="text" v-model="query.host"><br>
    File: <input type="text" v-model="query.file"><br>
    <button @click="listBackups(1)">run query</button>

    <h1>{{ count }} snapshots found</h1>
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
    <p>page <input type="number" v-model="page" @change="listBackups(page)"> / {{ pageLast }}</p>
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
#backup input[type=number] {
  width: 3em;
}
</style>

<script>
"use strict";

export default {
  data() {
    return {
      query: {
        host: "%",
        file: "%",
      },

      page: 1,
      pageRows: 1000,

      count: 0,
      snapshots: [],
    }
  },
  computed: {
    pageLast() {
      return Math.ceil(this.count / this.pageRows);
    },
  },
  created() {
    this.listBackups(1);
  },
  beforeDestroy() {
    clearInterval(this.timer);
  },
  methods: {
    async listBackups(page) {
      this.page = page;

      const fromWhereOrder = `FROM bookshelf_data_archive WHERE host LIKE '${this.query.host}' AND series LIKE '${this.query.file}' ORDER BY created DESC`;

      const cResp = await fetch(`/api/database/query`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(`SELECT COUNT(*) ${fromWhereOrder}`),
      });
      if(cResp.status != 200){
        return alert(await cResp.text());
      }
      this.count = parseInt((await cResp.json()).rows[0][0]);

      const resp = await fetch(`/api/database/query`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(`SELECT BINARY created, host, series, remote_id ${fromWhereOrder} LIMIT ${this.pageRows*(this.page-1)}, ${this.pageRows}`),
      });
      if(resp.status != 200){
        return alert(await resp.text());
      }

      this.snapshots = (await resp.json()).rows.map(([created, host, file, id]) => ({created, host, file, id}));
    }
  }
}
</script>
