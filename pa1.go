/*
“I Christopher Taliaferro (ch119541) affirm that this program is entirely my own work and that 
I have neither developed my code together with any another person, nor copied any code from any 
other person, nor permitted my code to be copied or otherwise used by any other person, nor 
have I copied, modified, or otherwise used programs created by others. I acknowledge that any 
violation of the above terms will be treated as academic dishonesty
*/

package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "strconv"
)

// Struct for each process
type process struct {
  arrival int
  burst int
  completionTime int
  finished bool
  initialBurst int
  identifier string
  processNumber int
  selected bool
  turnaround int
}

// Locate the index of a word from the input array
func getIndex(item string, list []string) (index int) {
  count := 0
  index = -1

  for i:= 0; i < len(list); i++ {
    if list[i] == item {
      index = count
      break
    }
    count++
  }
  return index
}

// Finds all processes and sorts by arrival time
func getProcesses(algorithm string, input string) ([]process) {
  var lines []string
  var processesArray []process
  count := 0

  // Read file
  file,_:= os.Open(input)
  defer file.Close()
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }

  // Extract processes info
  for i := 0; i < len(lines) - 1; i++ {
    // if statement : Skips first few lines to get to processes
    if (algorithm == "rr" && i > 3) || (algorithm != "rr" && i > 2) {
      currentProcess := new(process)
      var data []string
      data = strings.Split(lines[i], " ")

      // Fills processesArray with process structs
      if len(data) > 6 {
        currentProcess.processNumber = count
        currentProcess.identifier = data[2]
        currentProcess.arrival, _ = strconv.Atoi(data[4])
        currentProcess.initialBurst, _ = strconv.Atoi(data[6])
        currentProcess.burst, _ = strconv.Atoi(data[6])
        currentProcess.finished = false
      }
      processesArray = append(processesArray, *currentProcess)
      count++
    } 
  }

  // Bubble sort processesArray by arrival time
  for x := 0; x < len(processesArray); x++ {
    for y := 0; y < len(processesArray) - 1; y++ {
      if processesArray[y].arrival > processesArray[y + 1].arrival {
        processesArray[y], processesArray[y + 1] = processesArray[y + 1], processesArray[y]
      }
    }
  }
  return processesArray
}

// Calculate waiting time and turnaround time
func calculateTimes(processesArray []process, output *os.File)() {
  wait := 0
  // Calculate turnaround time
  for i := 0; i < len(processesArray); i++ {
    processesArray[i].turnaround = processesArray[i].completionTime - processesArray[i].arrival
  }

  // Output turnaround time & waiting time
  for i := 0; i < len(processesArray); i++ {
    for j := 0; j < len(processesArray); j++ {
      if (processesArray[j].processNumber == i) {
        wait = processesArray[j].turnaround - processesArray[j].initialBurst
        fmt.Fprintf(output,"%s wait%4d turnaround%4d\n", processesArray[j].identifier, wait , processesArray[j].turnaround)
      }
    }
  }
}

func firstComeFirstServe(proccessNum, runFor, quantum int, processesArray []process, output *os.File)() {
  running := false
  currentIndex := 0
  finish := 0
  start := 0 
  fmt.Fprintf(output, "%3d processes\n", proccessNum)
  fmt.Fprintln(output, "Using First-Come First-Served")

  for time := 0; time != runFor; time++ {
    // Process Arrived
    for i := 0; i < len(processesArray); i++{
      if time == processesArray[i].arrival {
        fmt.Fprintf(output, "Time%4d : %s arrived\n", time, processesArray[i].identifier)
        break
      }
    }

    // Process Finished
    if (time == start + processesArray[finish].burst) && running {
      fmt.Fprintf(output, "Time%4d : %s finished\n", time, processesArray[finish].identifier)
      running = false
      processesArray[finish].completionTime = time
      if currentIndex < len(processesArray) { currentIndex++ }
    }

    // Process Selection
    if !running && currentIndex < len(processesArray) && time >= processesArray[currentIndex].arrival { 
        fmt.Fprintf(output, "Time%4d : %s selected (burst%4d)\n", 
          time, processesArray[currentIndex].identifier, processesArray[currentIndex].burst)
        running = true
        finish = currentIndex
        start = time
    }

    // Idle
    if !running && processesArray[finish].arrival < time { fmt.Fprintf(output, "Time%4d : Idle\n", time) }

    // Algorithm Finished
    if time + 1 == runFor {
      fmt.Fprintln(output, "Finished at time ", time + 1)
      fmt.Fprintln(output)
    }
  }
  calculateTimes(processesArray, output)
}

func shortestJobFirst(proccessNum, runFor, quantum int, processesArray []process, output *os.File)() {
  running := false
  minBurstIndex := 0
  currentSetLimit := 0
  fmt.Fprintf(output, "%3d processes\n", proccessNum)
  fmt.Fprintln(output, "Using preemptive Shortest Job First")
  
  for time := 0; time != runFor; time++ {
    // Process Arrived
    for i := 0; i < len(processesArray); i++ {
      if time == processesArray[i].arrival {
        fmt.Fprintf(output, "Time%4d : %s arrived\n", time, processesArray[i].identifier)
        currentSetLimit = i
        break
      }
    }

    // Process Finished
    if processesArray[minBurstIndex].burst == 0 && running {
      fmt.Fprintf(output, "Time%4d : %s finished\n", time, processesArray[minBurstIndex].identifier)
      running = false
      processesArray[minBurstIndex].finished = true
      processesArray[minBurstIndex].completionTime = time
    }

    // Find Minimum Burst (shortest job)
    // If two processes have the same burst, pick whichever arrived earliest
    for j := 0; j <= currentSetLimit; j++ {
      if processesArray[j].finished == false && j != minBurstIndex{
        if processesArray[minBurstIndex].burst == 0 || 
           processesArray[j].burst < processesArray[minBurstIndex].burst ||
           (processesArray[j].burst == processesArray[minBurstIndex].burst && 
            processesArray[j].arrival < processesArray[minBurstIndex].arrival){
          processesArray[minBurstIndex].selected = false
          minBurstIndex = j
        } 
      }
    }

    // Process Selection
    if processesArray[minBurstIndex].burst != 0 && 
       processesArray[minBurstIndex].arrival <= time && 
       processesArray[minBurstIndex].selected == false { 

      fmt.Fprintf(output, "Time%4d : %s selected (burst%4d)\n", 
        time, processesArray[minBurstIndex].identifier, processesArray[minBurstIndex].burst)

      processesArray[minBurstIndex].selected = true
      running = true
    }

    // Idle
    if !running { fmt.Fprintf(output, "Time%4d : Idle\n", time) }

    // Running
    if processesArray[minBurstIndex].burst > 0 { processesArray[minBurstIndex].burst-- }

    // Algorithm Finished
    if time + 1 == runFor {
      fmt.Fprintln(output, "Finished at time ", time + 1)
      fmt.Fprintln(output)
    }
  }
  calculateTimes(processesArray, output)
}

func roundRobin(proccessNum, runFor, quantum int, processesArray []process, output *os.File)() {
  var queue []process
  running := false
  currentIndex := 0
  quantumLeft := 0
  fmt.Fprintf(output, "%3d processes\n", proccessNum)
  fmt.Fprintln(output, "Using Round-Robin")
  fmt.Fprintln(output, "Quantum  ", quantum)
  fmt.Fprintln(output)

  for time := 0; time != runFor; time++ {
    // Process Arrived
    for i := 0; i < len(processesArray); i++ {
      if time == processesArray[i].arrival {
        fmt.Fprintf(output, "Time%4d : %s arrived\n", time, processesArray[i].identifier)
        queue = append(queue, processesArray[i])
        break
      }
    }

    // Process Finished
    if processesArray[currentIndex].burst == 0 && 
       processesArray[currentIndex].finished == false && running {
      fmt.Fprintf(output, "Time%4d : %s finished\n", time, processesArray[currentIndex].identifier)
      running = false
      quantumLeft = 0
      queue = queue[1:]
      processesArray[currentIndex].finished = true
      processesArray[currentIndex].completionTime = time
    }
     
    // Process Selection
    if len(queue) > 0 && quantumLeft == 0 { 

      if len(queue) > 1 && running == true {
        queue = append(queue, queue[0])
        queue = queue[1:]
      }

      for i := 0; i < len(processesArray); i++ {
        if queue[0].identifier == processesArray[i].identifier {
          currentIndex = i
        }
      }

      fmt.Fprintf(output, "Time%4d : %s selected (burst%4d)\n", 
        time, processesArray[currentIndex].identifier, processesArray[currentIndex].burst)

      processesArray[currentIndex].selected = true
      running = true
      quantumLeft = quantum
    }

    // Idle
    if !running { fmt.Fprintf(output, "Time%4d : Idle\n", time) }

    // Running 
    if processesArray[currentIndex].burst > 0 {
      processesArray[currentIndex].burst--
      quantumLeft--
    }

    // Algorithm Finished
    if time + 1 == runFor {
      fmt.Fprintln(output, "Finished at time ", time + 1)
      fmt.Fprintln(output)
    }
  }
  calculateTimes(processesArray, output)
}

func main () {
  var data [] string
  
  // Set up input and output file
  file := os.Args[1]
  outputFile := os.Args[2]
  input,_:= os.Open(file)
  output,_:= os.Create(outputFile)
  defer input.Close()

  // Scanning each word
  scanner := bufio.NewScanner(input)
  scanner.Split(bufio.ScanWords)

  // data[] array will contain entire input file
  for scanner.Scan() { data = append(data, scanner.Text()) }

  // Parse input file and initialize attributes
  index := getIndex("processcount", data)
  processCount, _ := strconv.Atoi(data[index + 1])

  index = getIndex("runfor", data)
  runFor, _ := strconv.Atoi(data[index + 1])

  index = getIndex("quantum", data)
  quantum, _ := strconv.Atoi(data[index + 1])

  index = getIndex("use", data)
  algorithm := data[index + 1]

  // Gets processes and sorts them by arrival
  processes := getProcesses(algorithm, file)

  if algorithm == "fcfs" {
    firstComeFirstServe(processCount, runFor, quantum, processes, output)
  } else if algorithm == "sjf" {
    shortestJobFirst(processCount, runFor, quantum, processes, output)
  } else if algorithm == "rr" {
    roundRobin(processCount, runFor, quantum, processes, output)
  }

  defer output.Close()
}
