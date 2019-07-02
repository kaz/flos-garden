<template>
  <section id="cnc">
    <h1>Target</h1>
    <label :key="'tgt-'+node" v-for="node in nodes">
      <input type="checkbox" v-model="targets" :value="node"> {{ node }}
    </label>
    <p>
      <a href="javascript:" @click="selectAll(true)">check all</a>,
      <a href="javascript:" @click="selectAll(false)">uncheck all</a>
    </p>

    <h1>Command</h1>
    <textarea spellcheck="false" v-model="command"></textarea><br>

    <button @click="execute(false)">execute (oneshot)</button>
    period: <input type="number" v-model="period">
    <button @click="execute(true)" :disabled="timer">execute (periodically)</button>
    <button @click="stop()" :disabled="!timer">stop</button>

    <div :key="'out-'+target" v-for="target in targets">
      <h2>Result: {{ target }}</h2>
      <pre>{{ outputs[target] }}</pre>
    </div>

    <h1>Errors</h1>
    <pre class="err">{{ errors.join("\n") }}</pre>
  </section>
</template>

<style>
#cnc textarea {
  width: 60vw;
  height: 10em;
}
#cnc pre {
  padding: 1em;
  background-color: #eee;
}
#cnc pre.err {
  background-color: #fff;
  color: #f00;
}
</style>

<script>
"use strict";

export default {
  data() {
    return {
      nodes: [],
      targets: [],

      command: "",
      outputs: {},
      errors: [],

      period: 5,
      timer: null,
    }
  },
  created() {
    this.loadInstances();
  },
  beforeDestroy() {
    this.stop();
  },
  methods: {
    async loadInstances() {
      const resp = await fetch(`/api/collector/instance`);
      if(resp.status != 200){
        return alert(await resp.text());
      }
      this.nodes = (await resp.json()).sort();
    },
    selectAll(checked) {
      if(checked){
        this.targets = [].concat(this.nodes);
      }else{
        this.targets = [];
      }
    },
    stop() {
      clearInterval(this.timer);
      this.timer = null;
    },
    execute(interval) {
      this.outputs = {};
      this.errors = [];

      this.targets.forEach(target => {
        this.$set(this.outputs, target, "waiting ...");
      });

      if(interval){
        this.timer = setInterval(this._execute, this.period * 1000);
      }
      this._execute();
    },
    _execute() {
      this.targets.forEach(target => {
        this.executeOnNode(target);
      });
    },
    async executeOnNode(target) {
      const resp = await fetch(`/api/cnc/${target}/shell`, {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify(this.command),
      });
      if(resp.status != 200){
        this.errors.push(`[${target}] ${await resp.text()}`);
        return;
      }

      this.$set(this.outputs, target, await resp.json());
    }
  }
}
</script>
