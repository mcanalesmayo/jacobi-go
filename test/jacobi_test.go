package test

import (
	"fmt"
	"github.com/mcanalesmayo/jacobi-go"
	"github.com/mcanalesmayo/jacobi-go/model/matrix"
	"testing"
)

func TestRunJacobi(t *testing.T) {
	initialValue, nDim, maxIters, tolerance := 0.5, 16, 1000, 1.0e-4

	expectedMat := matrix.Matrix{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0.9959519906066812, 0.992048629498629, 0.9884275423329139, 0.985213226583245, 0.9825125033671933, 0.9804117347050112, 0.9789755934668327, 0.9782468926905672, 0.9782468926905674, 0.9789755934668328, 0.9804117347050113, 0.9825125033671933, 0.985213226583245, 0.988427542332914, 0.9920486294986288, 0.9959519906066812, 1},
		{1, 0.9917761439207942, 0.9838480345588239, 0.9764964757980613, 0.9699744961014689, 0.9644980616424679, 0.9602407391893988, 0.9573317425066765, 0.9558561828761315, 0.9558561828761315, 0.9573317425066767, 0.9602407391893988, 0.9644980616424679, 0.9699744961014689, 0.9764964757980611, 0.9838480345588239, 0.9917761439207942, 1},
		{1, 0.9873375848831656, 0.9751358327966154, 0.9638304717976247, 0.9538113368820703, 0.9454079748261823, 0.9388823496340637, 0.9344273707064148, 0.932168925678759, 0.932168925678759, 0.9344273707064148, 0.9388823496340635, 0.9454079748261823, 0.9538113368820702, 0.9638304717976247, 0.9751358327966152, 0.9873375848831656, 1},
		{1, 0.982486467369461, 0.9656218109029111, 0.9500160585551133, 0.9362087743245416, 0.9246490672809848, 0.9156876621675445, 0.9095782651290721, 0.9064839006071836, 0.9064839006071838, 0.9095782651290719, 0.9156876621675443, 0.9246490672809848, 0.9362087743245415, 0.9500160585551131, 0.9656218109029111, 0.9824864673694608, 1},
		{1, 0.9770479727705841, 0.9549697887586238, 0.9345793685961017, 0.9165841135035674, 0.9015589458194349, 0.8899405683597412, 0.8820360399152651, 0.8780377726987164, 0.8780377726987164, 0.8820360399152651, 0.8899405683597412, 0.9015589458194349, 0.9165841135035674, 0.9345793685961019, 0.9549697887586238, 0.9770479727705839, 1},
		{1, 0.9708083925438464, 0.9427730402573096, 0.9169559599677162, 0.894256122163736, 0.8753780174042682, 0.8608340747476796, 0.8509683982868644, 0.8459876487943995, 0.8459876487943996, 0.8509683982868644, 0.8608340747476794, 0.8753780174042682, 0.894256122163736, 0.9169559599677162, 0.9427730402573096, 0.9708083925438462, 1},
		{1, 0.9634940552605278, 0.9285182406485595, 0.8964487963621407, 0.8684052007392726, 0.8452168698967735, 0.8274463448853215, 0.8154424261975193, 0.8093984112551033, 0.8093984112551033, 0.8154424261975193, 0.8274463448853215, 0.8452168698967735, 0.8684052007392724, 0.8964487963621407, 0.9285182406485593, 0.9634940552605276, 1},
		{1, 0.9547370177171701, 0.9115289531126486, 0.8721662657440423, 0.8380195655532422, 0.8100176209381491, 0.7887179360512468, 0.7744141962476361, 0.7672389386315266, 0.7672389386315266, 0.7744141962476361, 0.7887179360512468, 0.8100176209381491, 0.8380195655532422, 0.8721662657440423, 0.9115289531126486, 0.9547370177171701, 1},
		{1, 0.944015424733358, 0.8908719356210277, 0.8429266306752394, 0.8018204769909252, 0.768508551796933, 0.7434337949553478, 0.7267304792941668, 0.718393878677383, 0.7183938786773831, 0.7267304792941669, 0.7434337949553478, 0.7685085517969329, 0.8018204769909251, 0.8429266306752394, 0.8908719356210277, 0.944015424733358, 1},
		{1, 0.9305429553247251, 0.8651940814198348, 0.8071062901440836, 0.7581579029840406, 0.7191540915108702, 0.69021768253323, 0.6711522441845498, 0.661701061952945, 0.661701061952945, 0.6711522441845498, 0.69021768253323, 0.7191540915108701, 0.7581579029840407, 0.8071062901440836, 0.8651940814198348, 0.9305429553247251, 1},
		{1, 0.9130493080851825, 0.8324261680109757, 0.7623957753738471, 0.7048697023648908, 0.6601100355475027, 0.6275543982740889, 0.6064151117699599, 0.5960284751530536, 0.5960284751530536, 0.6064151117699599, 0.6275543982740889, 0.6601100355475027, 0.7048697023648908, 0.7623957753738471, 0.8324261680109759, 0.9130493080851825, 1},
		{1, 0.8893089447799934, 0.7892244259335429, 0.7054125309221122, 0.6391114701307679, 0.5892130165415868, 0.5538685665122514, 0.5313484590075744, 0.5204072984000983, 0.5204072984000983, 0.5313484590075744, 0.5538685665122512, 0.5892130165415868, 0.6391114701307679, 0.7054125309221123, 0.7892244259335429, 0.8893089447799934, 1},
		{1, 0.8550340043951172, 0.7298915281070658, 0.6311246115780353, 0.5572144607282699, 0.5040745102785862, 0.46770895241979604, 0.44507952553134883, 0.43423490631555367, 0.43423490631555367, 0.44507952553134883, 0.46770895241979604, 0.5040745102785862, 0.5572144607282699, 0.6311246115780353, 0.7298915281070657, 0.8550340043951172, 1},
		{1, 0.8009962173331083, 0.6443023496549051, 0.5321537502532941, 0.45476969980912046, 0.40242511024470595, 0.36810878296901944, 0.3473433713177546, 0.33754667824704687, 0.33754667824704687, 0.34734337131775456, 0.36810878296901944, 0.40242511024470606, 0.4547696998091206, 0.5321537502532941, 0.6443023496549052, 0.8009962173331083, 1},
		{1, 0.7046958768692897, 0.5142610132484301, 0.3985540282523869, 0.32745912353093765, 0.28295313737802114, 0.2551884267707393, 0.23888640962963703, 0.2313184081396858, 0.2313184081396858, 0.23888640962963703, 0.2551884267707393, 0.28295313737802114, 0.32745912353093765, 0.39855402825238695, 0.5142610132484301, 0.7046958768692897, 1},
		{1, 0.5035587521074295, 0.30955564272988967, 0.2204352656312532, 0.1736786951141108, 0.1468809274968253, 0.1309635848361717, 0.12186542164679866, 0.11769811874997892, 0.11769811874997892, 0.12186542164679866, 0.1309635848361717, 0.1468809274968253, 0.1736786951141108, 0.2204352656312532, 0.3095556427298897, 0.5035587521074295, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	nThreadsExpectedSuccessCases := []int{1, 4}

	for _, nThreads := range nThreadsExpectedSuccessCases {
		fmt.Printf("Running simulation with matrix type='%s', initial value=%.4f, num dims=%d, max iterations=%d, tolerance=%.4f and num threads=%d\n",
			matrix.TwoDimMatrixType, initialValue, nDim, maxIters, tolerance, nThreads)

		actualMat, _, _ := jacobi.RunJacobi(matrix.TwoDimMatrixType, initialValue, nDim, maxIters, tolerance, nThreads)
		if !actualMat.CompareTo(expectedMat) {
			t.Errorf("Expected matrix doesn't match the actual one")
		}
	}
}
