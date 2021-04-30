/**
 * 直接console跑
 * 主观题最后两题自己勾
 */

prefix = 'item_'

// 客观题
radio_name_suffix = ['307', '308', '309', // 1
                     '311', '312', '313', '314', '315', // 2
                     '335', '336', '337', '338', '339', // 3
                     '341', '342', '343', '344', '345', // 4
                     '347', '348', '349', '350']; //5

for (let i = 0; i < radio_name_suffix.length; i++) {
    let radios = document.getElementsByName(prefix + radio_name_suffix[i]);
    let score = parseInt(Math.random() * 2 + 4);
    Array.from(radios).some( radio => {
        if(radio.value == score){
            radio.checked = true;
            return true;
        }
    });
}

// 主观题
// 1. 这门课程我最喜欢什么? 316
let ass_kiss = ['教师通过对课本的独到深入的讲解，达到了很好的教学效果，能结合多种教学手段，使学生对知识的掌握更深刻。',
        '教学内容重点突出，教学目的十分明确，教师具有极高的专业技能。',
        '授课方式新颖别致，激起同学们的兴趣，教师很注重互动，课堂学习氛围轻松愉快，真正达到了教学的目的要求。',
        '教师的教学效果极佳，可以使同学在领略知识魅力的同时提高自己实际技能。',
        '教师教课内容广大博深，高质量，高效率。教课内容新颖，独特，有个性。',
        '教师授课表现出来的激情和精神可以深深吸引并打动学生，希望我们的老师可以继续创新，造出更多的精品课。',
        '教师教学在书面浅显知识的基础上，进一步扩大了教学的知识的深度及广度，扩大了学生知识面，并且多方面培养学生的思考问题的能力。',
        '课上教师很注意与学生的互动环节，增强了课堂气氛，使教学效果更加显著。',
        '教师课堂上的整体教学效果非常好，教师在教学方面极认真负责，教师的基本知识技能过硬。',
        '课上教师很注意与学生的互动环节，语言也很生动、形象。',
        '教师并未忽视同学们的自己动手的锻炼、课堂互动效果极好。',
        '教师上课认真负责，专业基础极技能高深，非常注重学生的实际动手能力。',
        '上课语言幽默，互动适当，演示精准精彩。学生上课出勤率高，教学效果极其明显。'];

let text_1 = document.getElementById(prefix + '316');
idx = parseInt(Math.random() * ass_kiss.length);
if (text_1 != null){
    text_1.value = ass_kiss[idx];
}

// 2. 我认为本课程应从哪些方面需要进一步改进和提高 317
let text_2 = document.getElementById(prefix + "317");
let i = idx;
while(i == idx){
    i = parseInt(Math.random() * ass_kiss.length);
}
if (text_2 != null){
    text_2.value = ass_kiss[i] + '课程已经很不错了，希望老师以后能够有更多的进步' + '。';
}

// 3. 我平均每周在这门课程上花费多少小时? 318
let text_3 = document.getElementById(prefix + '318');
hour = parseInt(Math.random() * 4 + 2);
if (text_3 != null){
    text_3.value = '我平均每周在这门课上花费' + hour + '小时' + '。';
}

// 4. 在参与这门课之前，我对这个学科领域兴趣如何? 319
let items = ['就已经十分感兴趣', '有点兴趣但了解不多', '有一点兴趣','不是很感兴趣', '不是很感兴趣，但有听说过'];
idx = parseInt(Math.random() * items.length)
let text_4 = document.getElementById(prefix + '319');
if (text_4 != null){
    text_4.value = '在参与这门课之前，我对这个学科领域' + items[idx] + '。';
}

// 5. 我对该课程的课堂参与度（包括出勤、回答问题等） 320
// attendance = ['出勤完整','有课必到','出勤率较高, 偶尔请假'];
// a_idx = parseInt(Math.random() * attendance.length);
participant = ['表现积极', '表现较好', '表现积极', '表现不错','表现一般'];
p_idx = parseInt(Math.random() * participant.length);
let text_5 = document.getElementById(prefix + '320');
if (text_5 != null){
    text_5.value = '我出勤完整，' + '在课堂回答问题方面' + participant[p_idx] + '。';
}