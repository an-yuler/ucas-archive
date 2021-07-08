<template>
  <div class="q1">
    <div class="query-form">
      <el-alert
        title="任务一"
        :closable="false"
        type="info"
        description="查询两个公司之间的关联路径。例如输入公司“招商局轮船股份有限公司”和“招商银行股份有限公司”"
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
          Api.query1(vo).then((response) => {
            if (response.data.code != 100) {
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
    },
    onUseExample() {
      this.form.company1 = "招商局轮船股份有限公司";
      this.form.company2 = "招商银行股份有限公司";
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
</style>
