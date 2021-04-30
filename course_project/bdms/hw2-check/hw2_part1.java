/**
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Modified by Shimin Chen to demonstrate functionality for Homework 2
// April-May 2015

import java.io.BufferedReader;
import java.io.IOException;
import java.io.StringReader;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.DoubleWritable;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.util.GenericOptionsParser;

public class Hw2Part1{

    // This is the Mapper class
    // reference: http://hadoop.apache.org/docs/r2.6.0/api/org/apache/hadoop/mapreduce/Mapper.html
    //
    final private static String SPLIT_CHAR = " ";
    public static class CountMapper
            extends Mapper<Object, Text, Text, Text> {

        private Text keyout = new Text();
//        private Text valout = new Text();
        private Text valout = new Text();

        @Override
        public void map(Object key, Text value, Context context) throws IOException, InterruptedException {
            BufferedReader reader = new BufferedReader(new StringReader(value.toString()));
            String line = null;
            while((line = reader.readLine()) != null){
                String[] arr = line.split("\\s+");
                if(arr.length != 3) {
                    continue;
                }
                keyout.set(arr[0] + SPLIT_CHAR + arr[1]);
                valout.set(arr[2]);
                context.write(keyout, valout);
            }
        }
    }

    public static class CountCombiner
            extends Reducer<Text, Text, Text, Text> {
//        private DoubleWritable result = new DoubleWritable();
        private Text valout = new Text();
        @Override
        public void reduce(Text key, Iterable<Text> values,
                           Context context
        ) throws IOException, InterruptedException {
            int count = 0;
            double sum = 0;
            for (Text val : values) {
                count += 1;
                sum += Double.parseDouble(val.toString());
            }
            valout.set(count + SPLIT_CHAR + sum);
            context.write(key, valout);
        }
    }

    // This is the Reducer class
    // reference http://hadoop.apache.org/docs/r2.6.0/api/org/apache/hadoop/mapreduce/Reducer.html
    // We want to control the output format to look at the following: count of word = count
    public static class AverageReducer
            extends Reducer<Text, Text, Text, Text> {

        private Text resultKey = new Text();
        private Text resultValue = new Text();
        private byte[] blank;

        @Override
        protected void setup(Context context) {
            try {
                blank = Text.encode(" ").array();
            } catch (Exception e) {
                blank = new byte[0];
            }
        }

        @Override
        public void reduce(Text key, Iterable<Text> values,
                           Context context
        ) throws IOException, InterruptedException {
            int count = 0;
            double sum = 0;
            for (Text val : values) {
                String[] valueArr = val.toString().split(SPLIT_CHAR);
                count += Integer.parseInt(valueArr[0]);
                sum += Double.parseDouble(valueArr[1]);
            }
            // compute average
            String average = String.format("%.3f", sum / count);

            // generate result key
            String[] keyArr = key.toString().split(SPLIT_CHAR);
            resultKey.set(keyArr[0]);
            resultKey.append(blank, 0, blank.length);
            resultKey.append(keyArr[1].getBytes(), 0, keyArr[1].getBytes().length);

            // generate result value
            resultValue.set(String.valueOf(count).getBytes());
            resultValue.append(blank, 0, blank.length);
            resultValue.append(average.getBytes(), 0, average.getBytes().length);

            context.write(resultKey, resultValue);
        }
    }

    public static void main(String[] args) throws Exception {
        Configuration conf = new Configuration();
        String[] otherArgs = new GenericOptionsParser(conf, args).getRemainingArgs();
        if (otherArgs.length < 2) {
            System.err.println("Usage: wordcount <in> [<in>...] <out>");
            System.exit(2);
        }

        Job job = Job.getInstance(conf, "account from src to dst");

        job.setJarByClass(Hw2Part1.class);

        job.setMapperClass(CountMapper.class);
        job.setCombinerClass(CountCombiner.class);
        job.setReducerClass(AverageReducer.class);

        job.setMapOutputKeyClass(Text.class);
        job.setMapOutputValueClass(Text.class);

        job.setOutputKeyClass(Text.class);
        job.setOutputValueClass(Text.class);

        // add the input paths as given by command line
        for (int i = 0; i < otherArgs.length - 1; ++i) {
            FileInputFormat.addInputPath(job, new Path(otherArgs[i]));
        }

        // add the output path as given by the command line
        FileOutputFormat.setOutputPath(job,
                new Path(otherArgs[otherArgs.length - 1]));

        System.exit(job.waitForCompletion(true) ? 0 : 1);
    }
}

