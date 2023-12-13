<script setup>
import {reactive} from 'vue'
import {ScreenShot} from '../../wailsjs/go/main/App'
import jsQR from 'jsqr';
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const centerDialogVisible = ref(false)
const cropper = ref(null)
const data = reactive({
  img:"",
  qrcode:"",
})

function screenShot(){
  window.runtime.WindowHide()
  ScreenShot().then(result => {
    data.img = result
    window.runtime.WindowShow()
    centerDialogVisible.value = true
  })

}

function screenCut(){
  cropper.value.getCropData(result => {
    // do something
    data.qrcode = result
    console.log(result)
    decodeQr()
  })
}

function decodeQr(){
  
  const image = new Image();
  image.src = `${data.qrcode}`;

  image.onload = () => {
    const canvas = document.createElement('canvas');
    const context = canvas.getContext('2d');
    canvas.width = image.width;
    canvas.height = image.height;
    context.drawImage(image, 0, 0, image.width, image.height);
    const imageData = context.getImageData(0, 0, image.width, image.height);
    const code = jsQR(imageData.data, imageData.width, imageData.height);
    if (code) {
      // 在此处理识别到的二维码结果
      ElMessage({
        message: '识别到二维码内容：' + code.data,
        type: 'success',
      })
    } else {
      ElMessage({
        message: '未识别到二维码',
        type: 'warning',
      })
    }
}

}

</script>

<template>
  <main>

  <el-button text @click="screenShot">
    截图识别二维码
  </el-button>

  <el-dialog v-model="centerDialogVisible" title="请选取二维码" width="80%" center>
    <span>
      <vueCropper style="height:400px;width:100%"
      ref="cropper"
      :img="data.img"
      autoCrop
      centerBox
      outputSize="1"
      ></vueCropper>
    </span>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="centerDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="screenCut">
          确认选取
        </el-button>
      </span>
    </template>
  </el-dialog>


  


  </main>
</template>

<style scoped>
.dialog-footer button:first-child {
  margin-right: 10px;
}
</style>
