<template>
  <div class="q1">
    <div class="query-form">
      <el-alert
        title="任务三"
        :closable="false"
        type="info"
        description="环形持股查询，判断两家公司(5跳内)是否存在环形持股现象，环形持股是指两家公司彼此持有对方的股份。例如：输入“A”和“C”"
      >
      </el-alert>
      <el-form :rules="rules" ref="form" :model="form" label-width="80px">
        <el-form-item label="公司一" prop="company1">
          <el-col :span="7">
            <el-input
              v-model="form.company1"
              placeholder="请输入公司名"
            ></el-input>
          </el-col>
        </el-form-item>
        <el-form-item label="公司二" prop="company2">
          <el-col :span="7">
            <el-input
              v-model="form.company2"
              placeholder="请输入公司名"
            ></el-input>
          </el-col>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit('form')">查询</el-button>
          <el-button type="info" @click="onUseExample">示例</el-button>
          <el-button @click="onClearForm('form')">清空</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="res-show">
      <div class="alert-show">
        <el-alert
          v-show="isAlertSuccess"
          :closable="false"
          title="不存在循环持股"
          type="success"
          effect="dark"
        >
        </el-alert>
        <el-alert
          v-show="isAlertErr"
          title="存在循环持股"
          :closable="false"
          type="error"
          effect="dark"
        >
        </el-alert>
      </div>
      <el-table
        :data="responseVO.results"
        border
        height="700"
        style="width: 100%"
      >
        <el-table-column type="index" width="80" label="id"> </el-table-column>
        <!-- <el-table-column prop="chain" label="chains" width="580"> </el-table-column> -->
        <el-table-column prop="chain" label="chain">
          <template slot-scope="scope">
            {{ scope.row.chain.join("→") }}
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script>
import Api from "@/api/index.js";
export default {
  data() {
    return {
      form: {
        company1: "",
        company2: "",
      },
      responseVO: "",
      isAlertSuccess: false,
      isAlertErr: false,
      rules: {
        company1: [
          { required: true, message: "公司名一，必填项", trigger: "blur" },
        ],
        company2: [
          { required: true, message: "公司名二，必填项", trigger: "blur" },
        ],
      },
    };
  },
  methods: {
    onSubmit(formName) {
      this.$refs[formName].validate((valid) => {
        var vo = { params: [this.form.company1, this.form.company2] };
        if (valid) {
          Api.query3(vo).then((response) => {
            if (response.data.code == 100) {
              this.isAlertErr = true;
              this.isAlertSuccess = false;
            } else if (response.data.code == 301) {
              this.isAlertSuccess = true;
              this.isAlertErr = false;
            } else {
              this.isAlertSuccess = false;
              this.isAlertErr = false;
              this.$message({
                showClose: true,
                message: "出错了",
                type: "error",
              });
              return;
            }

            this.responseVO = response.data;
            this.$message({
              showClose: true,
              message: "查询完成",
              type: "success",
            });
          });
        } else {
          console.log("error submit!!");
          return false;
        }
      });
    },
    onClearForm(formName) {
      this.$refs[formName].resetFields();
      this.responseVO = "";
      this.isAlertSuccess = false;
      this.isAlertErr = false;
    },
    onUseExample() {
      this.form.company1 = "A";
      this.form.company2 = "C";
    },
  },
};
</script>
<style scoped>
.el-alert {
  margin: 10px 0;
}

.res-show {
  margin-top: 100px;
}
.alert-show {
  margin: 5px 0;
}
</style>
