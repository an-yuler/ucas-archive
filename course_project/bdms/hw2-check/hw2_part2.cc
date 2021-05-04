/**
 * @file PageRankVertex.cc
 * @author  Songjie Niu, Shimin Chen
 * @version 0.1
 *
 * @section LICENSE 
 * 
 * Copyright 2016 Shimin Chen (chensm@ict.ac.cn) and
 * Songjie Niu (niusongjie@ict.ac.cn)
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * 
 * @section DESCRIPTION
 * 
 * This file implements the PageRank algorithm using graphlite API.
 *
 */

#include <stdio.h>
#include <string.h>
#include <map>
#include <algorithm>
#include <math.h>

#include "GraphLite.h"
#define ll long long
#define VERTEX_CLASS_NAME(name) TriangleVertex##name

/* Message */
typedef int64_t VertexId;
enum State
{
    None,
    TO_IN,
    TO_OUT,
    FROM_IN,
    FROM_OUT,
};
/* MessageValue */
typedef struct
{
    VertexId src;
    State state; // 0 null, 1 in, 2 out
    VertexId data;
} Message;

/* VertexValue */
typedef struct
{
    ll in;
    ll out;
    ll through;
    ll cycle;
} TVNode;

/* List */
class List
{
private:
    vector<VertexId> in_list;
    vector<VertexId> out_list;

public:
    void addInList(VertexId id);
    void addOutList(VertexId id);
    bool searchInList(VertexId id);
    bool searchOutList(VertexId id);
    void print();
};
void List::addInList(VertexId id) { in_list.push_back(id); }

void List::addOutList(VertexId id) { out_list.push_back(id); }

bool List::searchInList(VertexId id)
{
    vector<VertexId>::iterator iter = find(in_list.begin(), in_list.end(), id);
    return iter != in_list.end();
}
bool List::searchOutList(VertexId id)
{
    vector<VertexId>::iterator iter = find(out_list.begin(), out_list.end(), id);
    return iter != out_list.end();
}

/* for debug */
void List::print()
{
    vector<VertexId>::iterator iter = in_list.begin();
    printf("in_list:[");
    while (iter != in_list.end())
    {
        if (iter == in_list.begin())
        {
            printf("%ld", *iter);
        }else{
            printf(", %ld", *iter);
        }
        
        iter++;
    }
    printf("]\n");
    iter = out_list.begin();
    printf("out_list[");
    while (iter != out_list.end())
    {
        if (iter == out_list.begin())
        {
            printf("%ld", *iter);
        }else{
            printf(", %ld", *iter);
        }
        iter++;
    }
    printf("]\n");
}

class VERTEX_CLASS_NAME(InputFormatter) : public InputFormatter
{
public:
    int64_t getVertexNum()
    {
        unsigned long long n;
        sscanf(m_ptotal_vertex_line, "%lld", &n);
        m_total_vertex = n;
        return m_total_vertex;
    }
    int64_t getEdgeNum()
    {
        unsigned long long n;
        sscanf(m_ptotal_edge_line, "%lld", &n);
        m_total_edge = n;
        return m_total_edge;
    }
    int getVertexValueSize()
    {
        m_n_value_size = sizeof(TVNode);
        return m_n_value_size;
    }
    int getEdgeValueSize()
    {
        m_e_value_size = sizeof(int);
        return m_e_value_size;
    }
    int getMessageValueSize()
    {
        m_m_value_size = sizeof(Message);
        return m_m_value_size;
    }
    void loadGraph()
    {
        unsigned long long last_vertex;
        unsigned long long from;
        unsigned long long to;
        double weight = 0;

        double value = 1;
        int outdegree = 0;

        const char *line = getEdgeLine();

        // Note: modify this if an edge weight is to be read
        //       modify the 'weight' variable

        sscanf(line, "%lld %lld", &from, &to);
        addEdge(from, to, &weight);

        last_vertex = from;
        ++outdegree;
        for (int64_t i = 1; i < m_total_edge; ++i)
        {
            line = getEdgeLine();

            // Note: modify this if an edge weight is to be read
            //       modify the 'weight' variable

            sscanf(line, "%lld %lld", &from, &to);
            if (last_vertex != from)
            {
                addVertex(last_vertex, &value, outdegree);
                last_vertex = from;
                outdegree = 1;
            }
            else
            {
                ++outdegree;
            }
            addEdge(from, to, &weight);
        }
        addVertex(last_vertex, &value, outdegree);
    }
};

class VERTEX_CLASS_NAME(OutputFormatter) : public OutputFormatter
{
public:
    void writeResult()
    {
        int64_t vid;
        TVNode value;
        char s[1024];

        ResultIterator r_iter;
        r_iter.getIdValue(vid, &value);

        int n = sprintf(s, "in: %lld\nout: %lld\nthrough: %lld\ncycle: %lld\n", value.in, value.out, value.through, value.cycle);
        writeNextResLine(s, n);
    }
};

// An aggregator that records a double value tom compute sum
class VERTEX_CLASS_NAME(Aggregator) : public Aggregator<ll>
{
public:
    void init()
    {
        m_global = 0;
        m_local = 0;
    }
    void *getGlobal()
    {
        return &m_global;
    }
    void setGlobal(const void *p)
    {
        m_global = *(ll *)p;
    }
    void *getLocal()
    {
        return &m_local;
    }
    void merge(const void *p)
    {
        m_global += *(ll *)p;
    }
    void accumulate(const void *p)
    {
        m_local += *(ll *)p;
    }
};

class VERTEX_CLASS_NAME() : public Vertex<TVNode, int, Message>
{
private:
    void info(vector<VertexId> &v)
    {
        /* print vector<> for debug */
        vector<VertexId>::iterator iter = v.begin();
        while (iter != v.end())
        {
            if (iter == v.begin())
            {
                printf("%ld", *iter);
            }
            else
            {
                printf(", %ld", *iter);
            }
            iter++;
        }
        printf("\n");
    }
    /* print map<> for debug */
    void info(map<VertexId, List> &m)
    {
        map<VertexId, List>::iterator iter = m.begin();
        while (iter != m.end())
        {
            printf("key: %ld\n", iter->first);
            iter->second.print();
            iter++;
        }
    }

public:
    void compute(MessageIterator *pmsgs)
    {
        /* super step 0: send empty message */
        if (getSuperstep() == 0)
        {
            Message msg = {getVertexId(), None, 0};
            sendMessageToAllNeighbors(msg);
        }
        /* super step 1: make adjacent vertex know about the in-vertexes and out-vertexes of current */
        else if (getSuperstep() == 1)
        {
            vector<VertexId> ins;
            /* pass in-list to to-vertex */
            for (; !pmsgs->done(); pmsgs->next())
            {
                Message msg = pmsgs->getValue();
                Message send = {getVertexId(), FROM_IN, msg.src};
                sendMessageToAllNeighbors(send);
                ins.push_back(msg.src);
            }
            /* pass in-list to from-vertex */
            for (int i = 0; i < ins.size(); i++)
            {
                for (int j = 0; j < ins.size(); j++)
                {
                    Message send = {getVertexId(), TO_IN, ins[j]};
                    sendMessageTo(ins[i], send);
                }
            }
            /* pass out-list to from-vertex and to-vertex */
            OutEdgeIterator outIter = getOutEdgeIterator();
            for (; !outIter.done(); outIter.next())
            {
                Message send = {getVertexId(), FROM_OUT, outIter.target()};
                sendMessageToAllNeighbors(send);
                for (int i = 0; i < ins.size(); i++)
                {
                    Message m = {getVertexId(), TO_OUT, outIter.target()};
                    sendMessageTo(ins[i], m);
                }
            }
        }
        /* supter step 3: compute in/out/through/cycle for each vertex, and accumulate the aggregator */
        else if (getSuperstep() == 2)
        {
            map<VertexId, List> frommap;
            map<VertexId, List> tomap;
            vector<VertexId> ins;
            vector<VertexId> outs;

            /* for debug */
            // printf("********************\n");
            for (; !pmsgs->done(); pmsgs->next())
            {
                Message msg = pmsgs->getValue();
                List list;
                if (msg.state == FROM_IN || msg.state == FROM_OUT)
                {
                    if (frommap.find(msg.src) == frommap.end())
                    {
                        frommap.insert(make_pair(msg.src, list));
                    }
                    list = frommap[msg.src];
                }
                else if (msg.state == TO_IN || msg.state == TO_OUT)
                {
                    if (tomap.find(msg.src) == tomap.end())
                    {
                        tomap.insert(make_pair(msg.src, list));
                    }
                    list = tomap[msg.src];
                }
                switch (msg.state)
                {
                case TO_IN:
                    list.addInList(msg.data);
                    tomap[msg.src] = list;
                    break;
                case FROM_IN:
                    list.addInList(msg.data);
                    frommap[msg.src] = list;
                    break;
                case TO_OUT:
                    list.addOutList(msg.data);
                    tomap[msg.src] = list;
                    break;
                case FROM_OUT:
                    list.addOutList(msg.data);
                    frommap[msg.src] = list;
                    break;
                }
            }
            /* for debug */
            // printf("VERTEX: %ld\n", getVertexId());
            // printf("frommap:\n");
            // info(frommap);
            // printf("tomap:\n");
            // info(tomap);
            // printf("********************\n");
            /* initialize ins and outs */
            map<VertexId, List>::iterator iter = frommap.begin();
            while (iter != frommap.end())
            {
                ins.push_back(iter->first);
                iter++;
            }
            iter = tomap.begin();
            while (iter != tomap.end())
            {
                outs.push_back(iter->first);
                iter++;
            }

            /* for debug */
            // printf("ins: \n");
            // info(ins);
            // printf("outs: \n");
            // info(outs);

            ll in_num, out_num, through_num, cycle_num;
            in_num = out_num = through_num = cycle_num = 0;
            /* compute in_num */
            for (int i = 0; i < ins.size(); i++)
            {
                for (int j = i + 1; j < ins.size(); j++)
                {
                    /* is connected ? */
                    in_num += frommap.at(ins[i]).searchInList(ins[j]) ? 1 : 0;
                    in_num += frommap.at(ins[i]).searchOutList(ins[j]) ? 1 : 0;

                }
            }

            /* compute out_num */
            for (int i = 0; i < outs.size(); i++)
            {
                for (int j = i + 1; j < outs.size(); j++)
                {
                    out_num += tomap.at(outs[i]).searchInList(outs[j]) ? 1 : 0;
                    out_num += tomap.at(outs[i]).searchOutList(outs[j]) ? 1 : 0;
                }
            }

            /* compute througn_num and cycle_num */
            for (int i = 0; i < ins.size(); i++)
            {
                for (int j = 0; j < outs.size(); j++)
                {
                    /* through_num */
                    through_num += frommap.at(ins[i]).searchOutList(outs[j]) ? 1 : 0;
                    /* cycle_num */
                    cycle_num += frommap.at(ins[i]).searchInList(outs[j]) ? 1 : 0;
                }
            }

            /* accumulate each aggregator */
            in_num += * (ll *)getAggrGlobal(0);
            out_num += * (ll *)getAggrGlobal(1);
            through_num += * (ll *)getAggrGlobal(2);
            cycle_num += * (ll *)getAggrGlobal(3);
        
            accumulateAggr(0, &in_num);
            accumulateAggr(1, &out_num);
            accumulateAggr(2, &through_num);
            accumulateAggr(3, &cycle_num);
        /* super step 4: it is last super step. we flush the result of aggregators to each vertex */
        }else{
            TVNode node;
            node.in += * (ll *)getAggrGlobal(0);
            node.out += * (ll *)getAggrGlobal(1);
            node.through += * (ll *)getAggrGlobal(2);
            node.cycle += * (ll *)getAggrGlobal(3);
            /* assign */
            *mutableValue() = node;
            /* finish computation */
            voteToHalt();
        }
    }
};

class VERTEX_CLASS_NAME(Graph) : public Graph
{
public:
    VERTEX_CLASS_NAME(Aggregator) * aggregator;

public:
    // argv[0]: PageRankVertex.so
    // argv[1]: <input path>
    // argv[2]: <output path>
    void init(int argc, char *argv[])
    {

        setNumHosts(5);
        setHost(0, "localhost", 1411);
        setHost(1, "localhost", 1421);
        setHost(2, "localhost", 1431);
        setHost(3, "localhost", 1441);
        setHost(4, "localhost", 1451);

        if (argc < 3)
        {
            printf("Usage: %s <input path> <output path>\n", argv[0]);
            exit(1);
        }

        m_pin_path = argv[1];
        m_pout_path = argv[2];

        aggregator = new VERTEX_CLASS_NAME(Aggregator)[4];
        regNumAggr(4);
        regAggr(0, &aggregator[0]);
        regAggr(1, &aggregator[1]);
        regAggr(2, &aggregator[2]);
        regAggr(3, &aggregator[3]);
    }

    void term()
    {
        delete[] aggregator;
    }
};

/* STOP: do not change the code below. */
extern "C" Graph *create_graph()
{
    Graph *pgraph = new VERTEX_CLASS_NAME(Graph);

    pgraph->m_pin_formatter = new VERTEX_CLASS_NAME(InputFormatter);
    pgraph->m_pout_formatter = new VERTEX_CLASS_NAME(OutputFormatter);
    pgraph->m_pver_base = new VERTEX_CLASS_NAME();

    return pgraph;
}

extern "C" void destroy_graph(Graph *pobject)
{
    delete (VERTEX_CLASS_NAME() *)(pobject->m_pver_base);
    delete (VERTEX_CLASS_NAME(OutputFormatter) *)(pobject->m_pout_formatter);
    delete (VERTEX_CLASS_NAME(InputFormatter) *)(pobject->m_pin_formatter);
    delete (VERTEX_CLASS_NAME(Graph) *)pobject;
}
