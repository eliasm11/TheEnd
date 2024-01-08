
/*
	"item": {
		"women": {
			"ts": {
				"1": {
					"xl": {
						"red": 1
					}
				},
				"2": {
					"sl": {
						"red": 1
					}
				}
			}
		}
	}
color: "black"
id: "p111"
size: "XL"
  */

let map = new Map()
for (const f in Order) {
	let o = Order[f]
	let container = f.split("/")[1]
	let kind = f.split("/")[2]
	if (typeof map[container] != "object") {
		map[container] = new Map()
	}
	if (typeof map[container][kind] !="object") {
		map[container][kind] = new Map()
	}
	for (const z of o) {

		if (typeof map[container][kind][z.id.replace("p", "")] != "object") {
			map[container][kind][z.id.replace("p", "")] = new Map()
		}

		if (typeof map[container][kind][z.id.replace("p", "")][z.size] != "object") {
			map[container][kind][z.id.replace("p", "")][z.size] = new Map()
		}

		if (typeof map[container][kind][z.id.replace("p", "")][z.size][z.color] != "object") {
		}

		map[container][kind][z.id.replace("p", "")][z.size][z.color] = z.qty



	}


}
console.log("the rnd")
console.log(JSON.stringify(map))



