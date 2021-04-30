import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataInputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.hbase.HBaseConfiguration;
import org.apache.hadoop.hbase.HColumnDescriptor;
import org.apache.hadoop.hbase.HTableDescriptor;
import org.apache.hadoop.hbase.TableName;
import org.apache.hadoop.hbase.client.HBaseAdmin;
import org.apache.hadoop.hbase.client.HTable;
import org.apache.hadoop.hbase.client.Put;


import java.io.*;
import java.math.BigDecimal;
import java.net.URI;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

/**
 * @author t
 */
public class Hw1Grp3 {
    public static void main(String[] args) throws IOException {
        if (args.length <= 0) {
            System.out.println("java Hw1Grp3 R=<file> groupby:R2 'res:count,avg(R3),max(R4)'");
            System.exit(1);
        }

        // main flow of sort based group-by
        Worker w = new Worker(args);
        w.readFromHdfs();
        w.compute();
        w.writeToHbase();

        // for local debugging
//        w.readFromLocalStorage();
//        w.compute();
//        w.printResult();
    }
}

class Worker{
    /**
     * There are three types of operations, including OP_COUNT, OP_AVG:[colId], OP_MAX:[colId]
     */
    private static final String OP_COUNT = "COUNT";
    private static final String OP_AVG = "AVG";
    private static final String OP_MAX = "MAX";

    /**
     * file path
     */
    private String filePath;
    /**
     * the index of group id, starts from 0
     */
    private int groupIdx;
    /**
     * general counter for each group, counting the number of items
     */
    private int counter;
    /**
     * original command line
     */
    private List<String> commandList;
    /**
     * the list contains all operations and its order is consistent with the input
     * specification: operationNum
     */
    private List<String> operationList;
    /**
     * each row of data after being parsed through {@link #parseEachLine(String)}.
     * specification: [rowNum * colNum]
     */
    private List<List<String>> itemList;
    /**
     * result list
     * specification: [groupNum * operationNum]
     */
    private List<List<String>> resultList;

    public Worker(String[] args) {
        counter = 0;
        commandList = new ArrayList<>();
        operationList = new ArrayList<>();
        itemList = new ArrayList<>();
        resultList = new ArrayList<>();
        parseCommand(args);
    }

    /**
     * parse command line args
     * initialize {@link #filePath},{@link #groupIdx} and {@link #operationList}
     * @param args command line args
     */
    private void parseCommand(String[] args) {
        filePath = args[0].substring(args[0].indexOf("=") + 1);
        groupIdx = Integer.parseInt(args[1].substring(args[1].lastIndexOf("R") + 1));
        String operation = args[2].substring(args[2].indexOf(":") + 1);
        String[] operations = operation.split(",");
        for (String op : operations) {
            commandList.add(op);
            if (op.contains("count")) {
                operationList.add(OP_COUNT);
            } else if (op.contains("avg")) {
                String id = op.substring(op.lastIndexOf('R') + 1, op.lastIndexOf(')'));
                operationList.add(OP_AVG + ':' + id);
            } else if (op.contains("max")) {
                String id = op.substring(op.indexOf('R') + 1, op.lastIndexOf(')'));
                operationList.add(OP_MAX + ':' + id);
            }
        }
    }

    /**
     * read file from local storage, for debug
     */
    public void readFromLocalStorage() throws IOException {
        BufferedReader reader = new BufferedReader(new FileReader(filePath));
        String s;
        while ((s = reader.readLine()) != null) {
            itemList.add(parseEachLine(s));
        }
        reader.close();
    }

    /**
     * read file from hdfs
     */
    public void readFromHdfs() throws IOException {
        String file = "hdfs://localhost:9000" + filePath;
        Configuration conf = new Configuration();
        FileSystem fs = FileSystem.get(URI.create(file), conf);
        FSDataInputStream in = fs.open(new Path(file));
        BufferedReader reader = new BufferedReader(new InputStreamReader(in));
        String s;
        while ((s = reader.readLine()) != null) {
            this.itemList.add(parseEachLine(s));
        }
        in.close();
        fs.close();
    }

    /**
     * split raw data by "|" to parse each line, then add put it into {@link #itemList}
     *
     * @param s raw data of each line in file
     * @return list containing the raw data that is split by "|"
     */
    private List<String> parseEachLine(String s) {
        String[] arr = s.split("\\|");
        return Arrays.stream(arr).collect(Collectors.toList());
    }

    /**
     * sort {@link #itemList} by the data at index {@link #groupIdx}
     */
    private void sort() {
        if (itemList.size() <= 1) {
            return;
        }
        itemList.sort(Comparator.comparing(o -> o.get(groupIdx)));
    }

    /**
     * compute for each group, and put result into {@link #resultList}
     */
    public void compute() {
        sort();

        String curGroup = itemList.get(0).get(groupIdx);
        List<Object> resultRound = new ArrayList<>(operationList.size());
        for (int i = 0; i < operationList.size(); i++) {
            resultRound.add(null);
        }

        for (int i = 0; i < itemList.size(); i++) {
            List<String> item = itemList.get(i);
            if (curGroup.equals(item.get(groupIdx))) {
                preprocess(item, resultRound);
            } else { // new group
                postprocess(curGroup, resultRound);

                // init for next group
                i--;
                // reset counter
                counter = 0;
                curGroup = item.get(groupIdx);
                resultRound = new ArrayList<>(operationList.size());
                for (int t = 0; t < operationList.size(); t++) {
                    resultRound.add(null);
                }
            }
        }
        // calculate for the last round
        postprocess(curGroup, resultRound);
    }

    /**
     * preprocess
     * the result of OP_COUNT and OP_MAX could be get immediately
     * but for OP_AVG, we must get the sum of num first
     * then compute the average value through {@link #postprocess(String, List)}
     *
     * @param item        row data for the specified group
     * @param resultRound intermediate result for the specified group
     */
    private void preprocess(List<String> item, List<Object> resultRound) {
        // preprocess for each line
        for (int j = 0; j < operationList.size(); j++) {
            String op = operationList.get(j);
            Object r = resultRound.get(j);
            if (op.equals(OP_COUNT)) {
                counter++;
                resultRound.set(j, r == null ? 1 : (int) r + 1);
            } else if (op.contains(OP_AVG)) {
                int index = Integer.parseInt(op.split(":")[1]);
                BigDecimal value = BigDecimal.valueOf(Double.parseDouble(item.get(index)));
                resultRound.set(j, r == null ? value : value.add((BigDecimal) r));
            } else if (op.contains(OP_MAX)) {
                int index = Integer.parseInt(op.split(":")[1]);
                int value = Double.valueOf(item.get(index)).intValue();
                resultRound.set(j, r == null ? value : Math.max((int) r, value));
            }
        }
    }

    /**
     * compute the average value for OP_AVG
     *
     * @param resultRound intermediate result for the specified group
     */
    private void postprocess(String group, List<Object> resultRound) {
        // calculate avg
        for (int j = 0; j < operationList.size(); j++) {
            if (operationList.get(j).contains(OP_AVG)) {
                    BigDecimal sum = (BigDecimal) resultRound.get(j);
                if (counter > 1){
                    resultRound.set(j, sum.divide(BigDecimal.valueOf(counter), 2, BigDecimal.ROUND_HALF_UP));
                }
                else if(counter == 1){
                    resultRound.set(j, sum.setScale(2,BigDecimal.ROUND_HALF_UP));
                }
            }
        }
        // add res for previous group
        resultRound.add(0, group);
        resultList.add(resultRound.stream().map(Object::toString).collect(Collectors.toList()));
    }

    public void writeToHbase() throws IOException {
        // create table if not exist
        final String TABLE_NAME = "Result";
        final String FAMILY_NAME = "res";
        Configuration conf = HBaseConfiguration.create();
        HBaseAdmin admin = new HBaseAdmin(conf);

        if (admin.tableExists(TABLE_NAME)) {
            System.out.println("Table already exists");
        } else {
            HTableDescriptor td = new HTableDescriptor(TableName.valueOf(TABLE_NAME));
            HColumnDescriptor cd = new HColumnDescriptor(FAMILY_NAME);
            td.addFamily(cd);
            admin.createTable(td);
            System.out.println("table " + TABLE_NAME + " created successfully");
        }
        admin.close();

        // add results into table
        HTable table = new HTable(conf, TABLE_NAME);
        for (List<String> result : resultList) {
            // row key
            Put put = new Put(result.get(0).getBytes());
            for (int i = 1; i < result.size(); i++) {
                put.add(FAMILY_NAME.getBytes(), commandList.get(i - 1).getBytes(), result.get(i).getBytes());
            }
            table.put(put);
        }
        table.close();
    }
    public void printResult(){
        resultList.forEach(System.out::println);
    }
}
